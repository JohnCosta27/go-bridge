package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
)

type NameType struct {
	Name     string
	Type     string
	Embedded bool
}

type Struct struct {
	Name   string
	Fields []NameType
}

type StructList []Struct

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

/*
 * Returned a topologically ordered list of structs,
 * This function is ugly and quite inefficient,
 * much room for improvement.
 */
func orderStructList(structList StructList) (StructList, error) {
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

func structAstToList(allAstStructs []StructWithName, astStructs []*ast.Field) ([]NameType, error) {
	structFields := make([]NameType, 0)

	for _, l := range astStructs {
		if len(l.Names) > 1 {
			return structFields, errors.New("More than one name returned")
		}

		fmt.Printf("AST TYPE: %T\n", l.Type)

		selectorExpr, ok := l.Type.(*ast.SelectorExpr)
		if ok {

			ident, ok := selectorExpr.X.(*ast.Ident)
			if !ok {
				return structFields, errors.New("Not a field access?")
			}

			//
			// when we get this, we should load the package into memory
			// (if it isn't already)
			// and get the type definitions from there.
			//

			fmt.Println(ident.Name, selectorExpr.Sel.Name)
		}

		fieldTypeIdent, ok := l.Type.(*ast.Ident)
		if !ok {
			return structFields, errors.New("Field type was more complicated")
		}

		fieldType := fieldTypeIdent.Name

		if len(l.Names) == 0 {
			// Embedded
			astFieldIndex := slices.IndexFunc(allAstStructs, func(ast StructWithName) bool {
				return ast.Name == fieldType
			})

			if astFieldIndex == -1 {
				return structFields, errors.New("Could not find embedded struct")
			}

			embeddedStruct := allAstStructs[astFieldIndex]

			embeddedStructFields, err := structAstToList(allAstStructs, embeddedStruct.Fields.List)
			if err != nil {
				return structFields, err
			}

			structFields = append(structFields, embeddedStructFields...)

			continue
		}

		fieldName := l.Names[0].Name

		if !ok {
			return structFields, errors.New("Field Type was more complicated, not supported yet")
		}

		nameType := NameType{Name: fieldName, Type: fieldType}
		structFields = append(structFields, nameType)
	}

	return structFields, nil
}

type StructWithName struct {
	*ast.StructType
	Name string
}

func getStructListFromAst(file *ast.File) (StructList, error) {
	structList := make(StructList, 0)

	astStructs := make([]StructWithName, 0)

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

			astStructs = append(astStructs, StructWithName{Name: typeSpec.Name.Name, StructType: structType})
		}

	}

	for _, s := range astStructs {
		structFields, err := structAstToList(astStructs, s.Fields.List)
		if err != nil {
			return structList, err
		}

		structList = append(structList, Struct{Name: s.Name, Fields: structFields})
	}

	return structList, nil
}

func ParseV2(entryFile string) (string, error) {
	fileDirectory := filepath.Dir(entryFile)

	//
	// First, let's read all the .go files from this DIR
	// as these can be used anywhere in entryFile
	// It does mean being heavier on resources initally.
	//

	files, err := os.ReadDir(fileDirectory)
	if err != nil {
		return "", err
	}

	goFiles := make([]fs.DirEntry, 0)

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		fileName := f.Name()

		if len(fileName) < 3 {
			continue
		}

		if fileName[len(fileName)-3:] != ".go" {
			continue
		}

		goFiles = append(goFiles, f)
	}

	goFilesContent := make([]string, len(goFiles))
	for i, f := range goFiles {
		content, err := os.ReadFile(fileDirectory + "/" + f.Name())
		if err != nil {
			return "", err
		}

		goFilesContent[i] = string(content)
	}

	return Parse(goFilesContent)
}

/*
 * Takes Golang code as input,
 * And outputs the correct parsing code
 * for Valibot.
 */
func Parse(goCode []string) (string, error) {
	astFile := make([]*ast.File, len(goCode))

	for i, code := range goCode {
		parsedFile, err := parser.ParseFile(token.NewFileSet(), "", code, 0)
		if err != nil {
			return "", err
		}

		astFile[i] = parsedFile
	}

	totalStructList := make(StructList, 0)

	for _, ast := range astFile {
		structList, err := getStructListFromAst(ast)
		if err != nil {
			return "", err
		}

		totalStructList = append(totalStructList, structList...)
	}

	newStructList, err := orderStructList(totalStructList)
	if err != nil {
		return "", err
	}

	valibotOutput, err := structsToValibot(newStructList)
	if err != nil {
		return "", err
	}

	return valibotOutput, nil
}
