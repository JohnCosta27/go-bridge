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

type FieldInfo struct {
	Name string
	Type string

	// Only one can be true at a time.
	// However, this isn't enforced on type level.
	// Consider adding an interface type and do type-checking.

	Embedded bool
	Array    bool
}

type Struct struct {
	Name   string
	Order  uint
	Fields []FieldInfo
}

type StructList []Struct

type OrderedStructType struct {
	*ast.StructType

	Order    uint
	FromFile string
	File     *ast.File
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

func (p *Parser) consumeFile(file *ast.File, fileName string) string {
	allStructs := make(NameToStructPos)
	packageName := file.Name.Name

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
				FromFile:   fileName,
				File:       file,
			}

			order++
		}
	}

	existingModuleStructs, exists := p.moduleStructs[packageName]
	if !exists {
		p.moduleStructs[packageName] = allStructs
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

		p.consumeFile(astFile, fileName)
	}

	return packageName, nil
}

func (p *Parser) parseEmbeddedStructField(orderedStruct OrderedStructType, structName string) ([]FieldInfo, error) {
	moduleStruct, exists := p.moduleStructs[orderedStruct.File.Name.Name]
	if !exists {
		return []FieldInfo{}, errors.New("Could not find package of embedded struct")
	}

	astF, exists := moduleStruct[structName]
	if !exists {
		return []FieldInfo{}, errors.New("Chould not find embedded struct")
	}

	embeddedStructFields, err := p.parseStruct(astF)
	if err != nil {
		return []FieldInfo{}, err
	}

	return embeddedStructFields, nil
}

func (p *Parser) parseDependencyField(orderedStruct OrderedStructType, fieldName string, expr *ast.SelectorExpr) ([]FieldInfo, error) {
	packageName, ok := expr.X.(*ast.Ident)
	if !ok {
		return []FieldInfo{}, errors.New("Could not parse dependency field")
	}

	structName := expr.Sel.Name

	depImport := orderedStruct.File.Imports[0].Path.Value
	depImport = depImport[len(p.projectPath)+2 : len(depImport)-1]

	p.consumeDir(depImport)

	packageStructs, exists := p.moduleStructs[packageName.Name]
	if !exists {
		return []FieldInfo{}, errors.New("Could not find package structs after consuming dir")
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
	p.moduleStructs[packageName.Name] = cleanPackageStructs

	return []FieldInfo{{Name: fieldName, Type: structName}}, nil
}

func (p *Parser) parseStructFieldType(orderedStruct OrderedStructType, fieldName string, field ast.Expr) ([]FieldInfo, error) {
	switch t := field.(type) {
	case *ast.Ident:
		return []FieldInfo{{Name: fieldName, Type: t.Name}}, nil
	case *ast.SelectorExpr:
		return p.parseDependencyField(orderedStruct, fieldName, t)
	case *ast.ArrayType:
		field, err := p.parseStructFieldType(orderedStruct, fieldName, t.Elt)
		if err != nil {
			return field, err
		}

		field[0].Array = true
		return field, err
	default:
		return []FieldInfo{}, errors.New(fmt.Sprintf("Currently, we don't support %T types.", field))
	}
}

func (p *Parser) parseStructField(orderedStruct OrderedStructType, field *ast.Field) ([]FieldInfo, error) {
	if len(field.Names) > 1 {
		return []FieldInfo{}, errors.New("More than one name returned")
	}

	if len(field.Names) == 0 {
		ident, ok := field.Type.(*ast.Ident)
		if !ok {
			return []FieldInfo{}, errors.New("Do not currently support non-ident types on embedded")
		}

		return p.parseEmbeddedStructField(orderedStruct, ident.Name)
	}

	return p.parseStructFieldType(orderedStruct, field.Names[0].Name, field.Type)
}

func (p *Parser) parseStruct(orderedStruct OrderedStructType) ([]FieldInfo, error) {
	structFields := make([]FieldInfo, 0)

	for _, field := range orderedStruct.Fields.List {
		processedFields, err := p.parseStructField(orderedStruct, field)
		if err != nil {
			return []FieldInfo{}, err
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
