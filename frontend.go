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

const EMBEDDED_DEP = "SUPER SECRET STRING FOR EMBEDDED DEPS"

type OrderedStructType struct {
	*ast.StructType
	File *ast.File

	StructName  string
	PackagePath string
	Order       uint
}

type NameToStructPos = map[string]OrderedStructType

// Map: ModuleName -> List of its structs
type ModuleStructs = map[string]NameToStructPos

// ====================== TODO ======================
// `moduleStructs` now is too complicated, because
// we are using the full paths to the packages as
// we know we can't have duplicated full paths.
//
// We could flatten this map,
// and make our code a bit nicer.
// ==================================================

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

			allStructs[packagePath+"-"+typeSpec.Name.Name] = OrderedStructType{
				StructType:  structType,
				StructName:  typeSpec.Name.Name,
				Order:       order,
				PackagePath: packagePath,

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

	astF, exists := moduleStruct[orderedStruct.PackagePath+"-"+structName]
	if !exists {
		return []StructField{}, errors.New("Could not find embedded struct")
	}

	embeddedStructFields, err := p.parseStruct(astF)
	if err != nil {
		return []StructField{}, err
	}

	return embeddedStructFields, nil
}

func (p *Parser) getPackagePath(imports []*ast.ImportSpec, packageName string) string {
	for _, i := range imports {
		v := i.Path.Value
		v = v[len(p.projectPath)+2 : len(v)-1]

		_, lastPath := filepath.Split(v)

		if lastPath == packageName {
			return v
		}
	}

	return ""
}

func (p *Parser) parseDependencyField(orderedStruct OrderedStructType, fieldName string, expr *ast.SelectorExpr) (StructField, error) {
	packageName, ok := expr.X.(*ast.Ident)
	if !ok {
		return BasicStructField{}, errors.New("Could not match type of package")
	}

	// ==== TODO: Refactor into seperate function ====

	depPackagePath := p.getPackagePath(orderedStruct.File.Imports, packageName.Name)
	structName := ""
	if depPackagePath != "" {
		structName = depPackagePath + "-" + expr.Sel.Name
	} else {
		structName = orderedStruct.PackagePath + "-" + expr.Sel.Name
	}

	// ====

	depImport := ""
	fullPath := ""

	for _, in := range orderedStruct.File.Imports {
		strippedValue := in.Path.Value[1 : len(in.Path.Value)-1]
		importPackageName := filepath.Base(strippedValue)

		if importPackageName == packageName.Name {
			depImport = importPackageName
			fullPath = strippedValue[len(p.projectPath)+1:]
			break
		}
	}

	if depImport == "" {
		return BasicStructField{}, errors.New("Could not find imported package")
	}

	p.consumeDir(fullPath)

	packageStructs, exists := p.moduleStructs[fullPath]
	if !exists {

		//
		// If the structs does not exist in a dependency, then this must be some
		// external dependency instead of a package dependency.
		//
		// We look for our own setup to see if we have a matching type,
		// otherwise we do with "any".
		//

		ident, ok := expr.X.(*ast.Ident)
		if !ok {
			return BasicStructField{}, errors.New("Expression type should be ident")
		}

		return UnknownStructField{FullType: ident.Name + "." + expr.Sel.Name, name: fieldName}, nil
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
	p.moduleStructs[fullPath] = cleanPackageStructs

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

	_, err := getJsType(valueIdent.Name)
	if err != nil {
		return MapStructField{name: fieldName, KeyType: keyIdent.Name, Value: BasicStructField{name: fieldName, Type: orderedStruct.PackagePath + "-" + valueIdent.Name}}, nil
	}

	return MapStructField{name: fieldName, KeyType: keyIdent.Name, Value: BasicStructField{name: fieldName, Type: valueIdent.Name}}, nil
}

func (p *Parser) parseStructFieldType(orderedStruct OrderedStructType, fieldName string, field ast.Expr) (StructField, error) {
	switch t := field.(type) {
	case *ast.Ident:
		// Same package dependant structs go in here.
		_, err := getJsType(t.Name)
		if err != nil {
			return BasicStructField{name: fieldName, Type: orderedStruct.PackagePath + "-" + t.Name}, nil
		}

		return BasicStructField{name: fieldName, Type: t.Name}, nil
	case *ast.SelectorExpr:
		return p.parseDependencyField(orderedStruct, fieldName, t)
	case *ast.StarExpr:
		return p.parseStructFieldType(orderedStruct, fieldName, t.X)
	case *ast.ArrayType:
		field, err := p.parseStructFieldType(orderedStruct, fieldName, t.Elt)
		if err != nil {
			return field, err
		}

		return ArrayStructField{name: field.Name(), Type: field}, err
	case *ast.MapType:
		return p.parseMapField(orderedStruct, fieldName, t)
	case *ast.StructType:
		orderedStruct.StructType = t
		fields, err := p.parseStruct(orderedStruct)

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
		switch field.Type.(type) {
		case *ast.Ident:
			ident := field.Type.(*ast.Ident)
			return p.parseEmbeddedStructField(orderedStruct, ident.Name)
		case *ast.SelectorExpr:
			selector := field.Type.(*ast.SelectorExpr)
			embeddedDepField, err := p.parseDependencyField(orderedStruct, EMBEDDED_DEP, selector)

			if err != nil {
				return []StructField{}, err
			}

			return []StructField{embeddedDepField}, nil
		default:
			return []StructField{}, errors.New(fmt.Sprintf("Do not currently support %T types on embedded", field.Type))
		}
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

	//
	// When we are processing embedded struct fields two things can happen.
	// 1. It's embedded from a struct in the same package. At which point we can
	//    grab it's fields very easily.
	//
	// 2. They are a dependency struct, which then are added as a BasicStructField,
	//    and the dependency is resolved later on.
	//
	// For approach 2, we need to do some post procesing, to get the struct fields
	// in the correct position.
	//
	// This creates two approaches for the same thing. Which is not ideal.
	// I actaully think we could move the embedded structs over to the backend,
	// but I'm not sure yet. Until then, we keep both approaches.
	//

	for _, s := range processedStructs {
		for _, field := range s.Fields {
			if field.Name() != EMBEDDED_DEP {
				continue
			}

			embeddedDepField := field.(BasicStructField)
			for _, ps := range processedStructs {
				if ps.Name != embeddedDepField.Type {
					continue
				}
			}

			fmt.Println(embeddedDepField.Type)
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
