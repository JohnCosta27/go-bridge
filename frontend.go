package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

type OrderedStructType struct {
	*ast.StructType
	File *ast.File

	Order uint
}

type NameToStructPos = map[string]OrderedStructType

// Map: ModuleName -> List of its structs
type ModuleStructs = map[string]NameToStructPos

type Parser struct {
	projectPath  string
	entryPackage string

	moduleStructs ModuleStructs
	outputStructs []Struct
}

func (p *Parser) consumeFile(file *ast.File, packagePath string) string {
	allStructs := make(NameToStructPos)

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

			allStructs[typeSpec.Name.Name] = OrderedStructType{
				StructType: structType,
				Order:      order,

				File: file,
			}

			order++
		}
	}

	existingModuleStructs, exists := p.moduleStructs[packagePath]
	if !exists {
		p.moduleStructs[packagePath] = allStructs
		return file.Name.Name
	}

	for k, v := range allStructs {
		existingModuleStructs[k] = v
	}

	return file.Name.Name
}

func (p *Parser) consumeDir(dirPath string) (string, error) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	packageName := ""

	for _, file := range files {
		fileName := file.Name()
		if len(fileName) < 3 || fileName[len(fileName)-3:] != ".go" {
			continue
		}

		fileContent, err := os.ReadFile(filepath.Join(dirPath, fileName))
		if err != nil {
			return "", err
		}

		astFile, err := parser.ParseFile(token.NewFileSet(), "", fileContent, 0)
		if err != nil {
			return "", err
		}

		packageName = astFile.Name.Name

		p.consumeFile(astFile, dirPath)
	}

	return packageName, nil
}

func (p *Parser) parseEmbeddedStructField(orderedStruct OrderedStructType, structName string) ([]StructField, error) {
	moduleStruct, exists := p.moduleStructs[orderedStruct.File.Name.Name]
	if !exists {
		return []StructField{}, errors.New("Could not find package of embedded struct")
	}

	astF, exists := moduleStruct[structName]
	if !exists {
		return []StructField{}, errors.New("Chould not find embedded struct")
	}

	embeddedStructFields, err := p.parseStruct(astF)
	if err != nil {
		return []StructField{}, err
	}

	return embeddedStructFields, nil
}

func (p *Parser) parseDependencyField(orderedStruct OrderedStructType, fieldName string, expr *ast.SelectorExpr) (StructField, error) {
	structName := expr.Sel.Name

	depImport := orderedStruct.File.Imports[0].Path.Value
	depImport = depImport[len(p.projectPath)+2 : len(depImport)-1]

	p.consumeDir(depImport)

	packageStructs, exists := p.moduleStructs[depImport]
	if !exists {
		return BasicStructField{}, errors.New("Could not find package structs after consuming dir")
	}

	//
	// Because we did `consumeDir`, as we don't know the exact file.
	// We should clean up, because we don't want all other structs that we
	// might not need present in our processing map, as they will make it
	// into the final output.
	//
	// This however, raises conserns over efficiency. We are deleting a dir
	// and perhaps re-reading it in the future. We can optimise this a tone.
	//

	cleanPackageStructs := make(NameToStructPos)
	cleanPackageStructs[structName] = packageStructs[structName]
	p.moduleStructs[depImport] = cleanPackageStructs

	return BasicStructField{name: fieldName, Type: structName}, nil
}

func (p *Parser) parseMapField(orderedStruct OrderedStructType, fieldName string, mapAst *ast.MapType) (StructField, error) {
	keyIdent, ok := mapAst.Key.(*ast.Ident)
	if !ok {
		return BasicStructField{}, errors.New(fmt.Sprintf("Only support %T as key of map type.", mapAst.Key))
	}

	valueIdent, ok := mapAst.Value.(*ast.Ident)
	if !ok {

		valueType, err := p.parseStructFieldType(orderedStruct, fieldName, mapAst.Value)
		if err != nil {
			return BasicStructField{}, err
		}

		return MapStructField{
			name:    fieldName,
			KeyType: keyIdent.Name,
			Value:   valueType,
		}, nil
	}

	return MapStructField{name: fieldName, KeyType: keyIdent.Name, Value: BasicStructField{name: fieldName, Type: valueIdent.Name}}, nil
}

func (p *Parser) parseStructFieldType(orderedStruct OrderedStructType, fieldName string, field ast.Expr) (StructField, error) {
	switch t := field.(type) {
	case *ast.Ident:
		return BasicStructField{name: fieldName, Type: t.Name}, nil
	case *ast.SelectorExpr:
		return p.parseDependencyField(orderedStruct, fieldName, t)
	case *ast.ArrayType:
		field, err := p.parseStructFieldType(orderedStruct, fieldName, t.Elt)
		if err != nil {
			return field, err
		}

		return ArrayStructField{name: field.Name(), Type: field}, err
	case *ast.MapType:
		return p.parseMapField(orderedStruct, fieldName, t)
	case *ast.StructType:
		fields, err := p.parseStruct(OrderedStructType{
			StructType: t,
		})

		if err != nil {
			return BasicStructField{}, err
		}

		return AnonStructField{name: fieldName, Fields: fields}, nil
	default:
		return BasicStructField{}, errors.New(fmt.Sprintf("Currently, we don't support %T types.", field))
	}
}

func (p *Parser) parseStructField(orderedStruct OrderedStructType, field *ast.Field) ([]StructField, error) {
	if len(field.Names) > 1 {
		return []StructField{}, errors.New("More than one name returned")
	}

	if len(field.Names) == 0 {
		ident, ok := field.Type.(*ast.Ident)
		if !ok {
			return []StructField{}, errors.New("Do not currently support non-ident types on embedded")
		}

		return p.parseEmbeddedStructField(orderedStruct, ident.Name)
	}

	structFields, err := p.parseStructFieldType(orderedStruct, field.Names[0].Name, field.Type)
	if err != nil {
		return []StructField{}, err
	}

	return []StructField{structFields}, nil
}

func (p *Parser) parseStruct(orderedStruct OrderedStructType) ([]StructField, error) {
	structFields := make([]StructField, 0)

	for _, field := range orderedStruct.Fields.List {
		processedFields, err := p.parseStructField(orderedStruct, field)
		if err != nil {
			return []StructField{}, err
		}

		structFields = append(structFields, processedFields...)
	}

	return structFields, nil
}

func (p *Parser) Parse() ([]Struct, error) {
	processedStructs := make([]Struct, 0)

	for len(p.moduleStructs) > 0 {

		for packageName, packageStructs := range p.moduleStructs {
			for structName, s := range packageStructs {
				fields, err := p.parseStruct(s)
				if err != nil {
					return []Struct{}, err
				}

				parsedStruct := Struct{
					Name:   structName,
					Order:  s.Order,
					Fields: fields,
				}

				processedStructs = append(processedStructs, parsedStruct)
			}

			delete(p.moduleStructs, packageName)
		}

	}

	return processedStructs, nil
}

func ParserFactory(entryFile string, givenProjectPath string) (Parser, error) {
	p := Parser{
		projectPath:   givenProjectPath,
		moduleStructs: make(ModuleStructs),
	}

	path := filepath.Dir(entryFile)

	mainPackage, err := p.consumeDir(path)
	if err != nil {
		return Parser{}, err
	}

	p.entryPackage = mainPackage
	return p, nil
}
