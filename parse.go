package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
)

type NameType struct {
	Name string
	Type string
}

type StructList map[string][]NameType

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

	return "", errors.New("Cannot find corresponding JS type")
}

func structsToValibot(structList StructList) (string, error) {
	valibotOutput := ""

	for structName, structType := range structList {
		localValidbotOutput := "const " + structName + " = object({\n"

		for _, fieldType := range structType {
			jsType, err := getJsType(fieldType.Type)
			if err != nil {
				return "", err
			}

			localValidbotOutput += "  " + fieldType.Name + ": " + jsType + "(),\n"
		}

		localValidbotOutput += "});"
		valibotOutput += "\n" + localValidbotOutput + "\n"
	}

	return valibotOutput, nil
}

func structAstToList(astStructs []*ast.Field) ([]NameType, error) {
	structFields := make([]NameType, 0)

	for _, l := range astStructs {
		if len(l.Names) != 1 {
			return structFields, errors.New("More than one name returned")
		}

		fieldName := l.Names[0].Name

		fieldTypeIdent, ok := l.Type.(*ast.Ident)
		fieldType := fieldTypeIdent.Name

		if !ok {
			return structFields, errors.New("Field type was more complicated")
		}

		if !ok {
			return structFields, errors.New("Field Type was more complicated, not supported yet")
		}

		structFields = append(structFields, NameType{Name: fieldName, Type: fieldType})
	}

	return structFields, nil
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

	structList := make(StructList)

	for _, dec := range parsedFile.Decls {
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

			structName := typeSpec.Name.Name

			structFields, err := structAstToList(structType.Fields.List)
			if err != nil {
				return "", err
			}

			structList[structName] = structFields
		}
	}

	valibotOutput, err := structsToValibot(structList)
	if err != nil {
		return "", err
	}

	return valibotOutput, nil
}