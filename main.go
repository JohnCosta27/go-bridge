package main

import (
	"go/parser"
	"go/token"
	"path/filepath"
)

func MainParse(entryFile string, givenProjectPath string) (string, error) {
	parser, err := ParserFactory(entryFile, givenProjectPath)
	if err != nil {
		return "", err
	}

	entryDir := filepath.Dir(entryFile)
	parser.consumeDir(entryDir)

	structs, err := parser.Parse()
	if err != nil {
		return "", err
	}

	structs, err = orderStructList(structs)
	if err != nil {
		return "", err
	}

	valibotOutput, err := structsToValibot(structs)
	if err != nil {
		return "", err
	}

	return valibotOutput, nil
}

func CodeParse(content string) (string, error) {
	p := Parser{
		moduleStructs: make(ModuleStructs),
	}

	astFile, err := parser.ParseFile(token.NewFileSet(), "", content, 0)
	if err != nil {
		return "", err
	}

	p.entryPackage = astFile.Name.Name
	p.consumeFile(astFile, astFile.Name.Name+".go")

	structs, err := p.Parse()
	if err != nil {
		return "", err
	}

	structs, err = orderStructList(structs)
	if err != nil {
		return "", err
	}

	valibotOutput, err := structsToValibot(structs)
	if err != nil {
		return "", err
	}

	return valibotOutput, nil
}
