package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"slices"
	"sort"
)

type FieldInfo struct {
	Name     string
	Type     string
	Embedded bool
}

type Struct struct {
	Name   string
	Order  uint
	Fields []FieldInfo
}

type StructList []Struct

type OrderedStructType struct {
	*ast.StructType
	Order uint
}

type NameToStructPos = map[string]OrderedStructType

var NoJsType = errors.New("Cannot find corresponding JS type")

func getJsType(goType string) (string, error) {
	switch goType {
	case "int":
		fallthrough
	case "int8":
		fallthrough
	case "int16":
		fallthrough
	case "int32":
		fallthrough
	case "int64":
		fallthrough
	case "float32":
		fallthrough
	case "float64":
		return "number", nil
	case "string":
		return "string", nil
	case "bool":
		return "boolean", nil
	}

	return "", NoJsType
}

// ==================================================
// Backend Methods.
//
// Responsible for ordering output, and producing,
// the resultant code.
// ==================================================

/*
 * Returned a topologically ordered list of structs,
 * This function is ugly and quite inefficient,
 * much room for improvement.
 */
func orderStructList(structList StructList) (StructList, error) {

	//
	// Because we used a map while getting our structs,
	// We need to order the structs based on order they came to us.
	// So we can get deterministic output.
	//
	sort.Slice(structList, func(i, j int) bool {
		return structList[i].Order < structList[j].Order
	})

	nodeMap := make([]*Node, 0)

	for _, s := range structList {
		node := &Node{Name: s.Name, Edges: make([]*Node, 0)}
		nodeMap = append(nodeMap, node)
	}

	nodeList := make([]*Node, 0)

	for i, node := range nodeMap {
		for _, field := range structList[i].Fields {
			_, err := getJsType(field.Type)

			if err != NoJsType {
				continue
			}

			nodeIndex := slices.IndexFunc(nodeMap, func(n *Node) bool {
				return n.Name == field.Type
			})

			node.Edges = append(node.Edges, nodeMap[nodeIndex])
		}

		nodeList = append(nodeList, node)
	}

	orderedList := make(StructList, len(structList))

	ordering := topologicalSort(nodeList)

	// Very inefficient linear search
	// TODO: Make it better.

	for orderIndex, n := range ordering {

		index := -1
		for i, s := range structList {
			if s.Name == n {
				index = i
				break
			}
		}

		if index == -1 {
			return structList, errors.New("Could not find index of node")
		}

		orderedList[orderIndex] = structList[index]
	}

	return orderedList, nil
}

func structsToValibot(structList StructList) (string, error) {
	valibotOutput := ""

	importedValidators := make([]string, 0)
	importedValidators = append(importedValidators, "object")

	for _, s := range structList {
		localValidbotOutput := "const " + s.Name + " = object({\n"

		for _, fieldType := range s.Fields {
			jsType, err := getJsType(fieldType.Type)

			if err == NoJsType {
				// Here, we must have a nested struct.
				localValidbotOutput += "  " + fieldType.Name + ": " + fieldType.Type + ",\n"

				continue
			}

			exist := slices.Index(importedValidators, jsType) != -1
			if !exist {
				importedValidators = append(importedValidators, jsType)
			}
			localValidbotOutput += "  " + fieldType.Name + ": " + jsType + "(),\n"

		}

		localValidbotOutput += "});"
		valibotOutput += "\n" + localValidbotOutput + "\n"
	}

	importLine := ""
	for _, validator := range importedValidators {
		importLine += validator + ", "
	}

	importLine = "import { " + importLine[:len(importLine)-2] + " } from 'valibot';\n"

	return "\n" + importLine + valibotOutput, nil
}

// ==================================================
// Frontend methods
//
// Responsible for parsing structs from AST, tracking
// dependencies and getting structs to usable format.
// ==================================================

func getEmbeddedStructFields(allAstStructs *NameToStructPos, structName string) ([]FieldInfo, error) {
	astF, exists := (*allAstStructs)[structName]
	if !exists {
		return []FieldInfo{}, errors.New("Chould not find embedded struct")
	}

	embeddedStructFields, err := structAstToList(allAstStructs, astF.Fields.List)
	if err != nil {
		return []FieldInfo{}, err
	}

	return embeddedStructFields, nil
}

func getPackageStructField(allAstStructs *NameToStructPos, expr *ast.SelectorExpr) error {
	fmt.Printf("%T\n", expr.X)

	return nil
}

func getSingleStructField(allAstStructs *NameToStructPos, field *ast.Field) ([]FieldInfo, error) {
	if len(field.Names) > 1 {
		return make([]FieldInfo, 0), errors.New("More than one name returned")
	}

	selectorExpr, ok := field.Type.(*ast.SelectorExpr)
	if ok {
		//
		// when we get this, we should load the package into memory
		// (if it isn't already)
		// and get the type definitions from there.
		//

		getPackageStructField(allAstStructs, selectorExpr)

		return []FieldInfo{}, errors.New("Not implemented yet.")
	}

	fieldTypeIdent, ok := field.Type.(*ast.Ident)
	if !ok {
		return []FieldInfo{}, errors.New("Field type was more complicated")
	}

	fieldType := fieldTypeIdent.Name
	if len(field.Names) == 0 {
		nestedStructFields, err := getEmbeddedStructFields(allAstStructs, fieldType)
		if err != nil {
			return []FieldInfo{}, err
		}

		return nestedStructFields, nil
	}

	fieldName := field.Names[0].Name
	return []FieldInfo{{Name: fieldName, Type: fieldType}}, nil
}

func structAstToList(allAstStructs *NameToStructPos, astStructs []*ast.Field) ([]FieldInfo, error) {
	structFields := make([]FieldInfo, 0)

	for _, field := range astStructs {
		processedField, err := getSingleStructField(allAstStructs, field)
		if err != nil {
			return structFields, err
		}

		structFields = append(structFields, processedField...)

	}

	return structFields, nil
}

func getStructListFromAst(file *ast.File) (StructList, error) {
	structList := make(StructList, 0)
	astStructs := make(NameToStructPos)

	var order uint = 0

	for _, dec := range file.Decls {
		typeDec, ok := dec.(*ast.GenDecl)

		if !ok {
			continue
		}

		for _, spec := range typeDec.Specs {

			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			astStructs[typeSpec.Name.Name] = OrderedStructType{StructType: structType, Order: order}

			order++
		}

	}

	for sName, s := range astStructs {
		structFields, err := structAstToList(&astStructs, s.Fields.List)
		if err != nil {
			return structList, err
		}

		structList = append(structList, Struct{Name: sName, Fields: structFields, Order: s.Order})
	}

	return structList, nil
}

func ParseV2(entryFile string) (string, error) {
	goContent, err := os.ReadFile(entryFile)
	if err != nil {
		return "", err
	}

	return Parse(string(goContent))
}

/*
 * Takes Golang code as input,
 * And outputs the correct parsing code
 * for Valibot.
 */
func Parse(goCode string) (string, error) {
	parsedFile, err := parser.ParseFile(token.NewFileSet(), "", goCode, 0)
	if err != nil {
		return "", err
	}

	structList, err := getStructListFromAst(parsedFile)
	if err != nil {
		return "", err
	}

	newStructList, err := orderStructList(structList)
	if err != nil {
		return "", err
	}

	valibotOutput, err := structsToValibot(newStructList)
	if err != nil {
		return "", err
	}

	return valibotOutput, nil
}
