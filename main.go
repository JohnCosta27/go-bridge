package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
)

func MainParse(entryFile string, givenProjectPath string) (string, error) {
	parser, err := ParserFactory(entryFile, givenProjectPath)
	if err != nil {
		return "", err
	}

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
	p.consumeFile(astFile, astFile.Name.Name)

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

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Please type an entry file")
		return
	}

	entryFile := args[0]
	output, err := MainParse(entryFile, "johncosta.tech/struct-to-types")
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
	}

	fmt.Print(output)
}
