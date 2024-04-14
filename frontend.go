package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
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
}

func (p *Parser) consumeFile(file *ast.File, fileName string) {
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
		return
	}

	for k, v := range allStructs {
		existingModuleStructs[k] = v
	}
}

func (p *Parser) consumeDir(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		fileName := file.Name()
		if len(fileName) < 3 || fileName[len(fileName)-3:] != ".go" {
			continue
		}

		fileContent, err := os.ReadFile(filepath.Join(dirPath, fileName))
		if err != nil {
			return err
		}

		astFile, err := parser.ParseFile(token.NewFileSet(), "", fileContent, 0)
		if err != nil {
			return err
		}

		p.entryPackage = astFile.Name.Name

		p.consumeFile(astFile, fileName)
	}

	return nil
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

	embeddedStructFields, err := p.parseStructs(astF)
	if err != nil {
		return []FieldInfo{}, err
	}

	return embeddedStructFields, nil
}

func (p *Parser) parseStructField(orderedStruct OrderedStructType, field *ast.Field) ([]FieldInfo, error) {
	if len(field.Names) > 1 {
		return []FieldInfo{}, errors.New("More than one name returned")
	}

	_, ok := field.Type.(*ast.SelectorExpr)
	if ok {
		//
		// when we get this, we should load the package into memory
		// (if it isn't already)
		// and get the type definitions from there.
		//

		// getPackageStructField(allAstStructs, s, selectorExpr)

		return []FieldInfo{}, errors.New("Not implemented yet.")
	}

	fieldTypeIdent, ok := field.Type.(*ast.Ident)
	if !ok {
		return []FieldInfo{}, errors.New("Field type was more complicated")
	}

	fieldType := fieldTypeIdent.Name
	if len(field.Names) == 0 {
		nestedStructFields, err := p.parseEmbeddedStructField(orderedStruct, fieldType)
		if err != nil {
			return []FieldInfo{}, err
		}

		return nestedStructFields, nil
	}

	fieldName := field.Names[0].Name
	return []FieldInfo{{Name: fieldName, Type: fieldType}}, nil
}

func (p *Parser) parseStructs(orderedStruct OrderedStructType) ([]FieldInfo, error) {
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
	mainPackage, exists := p.moduleStructs[p.entryPackage]
	if !exists {
		return []Struct{}, errors.New("Could not find any packages in entry packages")
	}

	processedStructs := make([]Struct, len(mainPackage))
	i := 0

	for structName, s := range mainPackage {
		fields, err := p.parseStructs(s)
		if err != nil {
			return []Struct{}, err
		}

		parsedStruct := Struct{
			Name:   structName,
			Order:  s.Order,
			Fields: fields,
		}

		processedStructs[i] = parsedStruct
		i++
	}

	return processedStructs, nil
}

func ParserFactory(entryFile string, givenProjectPath string) (Parser, error) {
	parser := Parser{
		projectPath:   givenProjectPath,
		moduleStructs: make(ModuleStructs),
	}

	mainDir := filepath.Dir(entryFile)
	err := parser.consumeDir(mainDir)
	if err != nil {
		return Parser{}, err
	}

	return parser, nil
}
