package main

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
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

func readProjectPath(goModPath string) (string, error) {
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(string(content), "\n") {
		if !strings.Contains(line, "module ") {
			continue
		}

		return strings.Split(line, " ")[1], nil
	}

	return "", errors.New("Could not find line containing module in go.mod file")
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Println("Please type an entry file")
		return
	}

	entryFile := args[0]

	packageDir, _ := filepath.Split(entryFile)
	withGoMod := packageDir + "go.mod"

	projectPath, err := readProjectPath(withGoMod)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		return
	}

	output, err := MainParse(entryFile, projectPath)
	if err != nil {
		fmt.Fprintln(os.Stdout, err)
		return
	}

	fmt.Print(output)
}
