package main

import (
	"errors"
	"slices"
	"sort"
	"strings"
)

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
	case "uint":
		fallthrough
	case "uint8":
		fallthrough
	case "uint16":
		fallthrough
	case "uint32":
		fallthrough
	case "uint64":
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

func getSpaces(indent uint) string {
	output := ""
	var i uint = 0
	for i = 0; i < indent; i++ {
		output += "  "
	}

	return output
}

// ==================================================
// Backend Methods.
//
// Responsible for ordering output, and producing,
// the resultant code.
// ==================================================

func recGetLastType(field StructField) string {
	switch t := field.(type) {
	case BasicStructField:
		return t.Type
	case MapStructField:
		return recGetLastType(t.Value)
	case ArrayStructField:
		return recGetLastType(t.Type)
	case AnonStructField:
		return recGetLastType(t.Fields[0])
	default:
		panic("Switch should be exhaustive")
	}
}

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

	//
	// At this point we could have duplicate names across packages.
	//

	nodeMap := make([]*Node, 0)

	for _, s := range structList {
		node := &Node{Name: s.Name, Edges: make([]*Node, 0)}
		nodeMap = append(nodeMap, node)
	}

	nodeList := make([]*Node, 0)

	for i, node := range nodeMap {
		for _, field := range structList[i].Fields {
			lastType := recGetLastType(field)

			_, err := getJsType(lastType)

			if err != NoJsType {
				continue
			}

			nodeIndex := slices.IndexFunc(nodeMap, func(n *Node) bool {
				return n.Name == lastType
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

// TODO: you are literally a set mate (with insert order)
func maybeAdd(validators map[string]uint, counter *uint, field string) {
	_, exists := validators[field]
	if exists {
		return
	}

	validators[field] = *counter
	*counter++
}

func getJsTypeOrType(t string) (bool, string) {
	jsType, err := getJsType(t)
	if err != nil {
		return true, t
	}

	return false, jsType
}

func appendedType(t string) string {
	isJsType, t := getJsTypeOrType(t)
	if isJsType {
		return t + "),\n"
	}

	return t + "()),\n"
}

func getStructFieldType(validators map[string]uint, nameMap map[string]string, counter *uint, field StructField, indent uint) (string, error) {
	switch t := field.(type) {
	case BasicStructField:
		jsType, err := getJsType(t.Type)

		if err == NoJsType {
			return nameMap[t.Type], nil
		}

		maybeAdd(validators, counter, jsType)
		return jsType + "()", nil
	case ArrayStructField:
		maybeAdd(validators, counter, "array")
		recValue, err := getStructFieldType(validators, nameMap, counter, t.Type, indent+1)
		if err != nil {
			return "", err
		}

		return "array(" + recValue + ")", nil
	case MapStructField:
		maybeAdd(validators, counter, "record")
		recValue, err := getStructFieldType(validators, nameMap, counter, t.Value, indent+1)
		if err != nil {
			return "", err
		}

		return "record(" + recValue + ")", nil
	case AnonStructField:
		output := "object({\n"
		for _, v := range t.Fields {
			fieldOutput, err := getSingleField(validators, nameMap, counter, v, indent+1)
			if err != nil {
				return "", err
			}

			output += fieldOutput
		}
		output += getSpaces(indent+1) + "})"
		return output, nil
	default:
		return "", errors.New("not implemented")
	}
}

func getSingleField(validators map[string]uint, nameMap map[string]string, counter *uint, field StructField, indent uint) (string, error) {
	typeValue, err := getStructFieldType(validators, nameMap, counter, field, indent)
	if err != nil {
		return "", err
	}

	return getSpaces(indent+1) + field.Name() + ": " + typeValue + ",\n", nil
}

func getName(namespacedName string) string {
	return strings.Split(namespacedName, "-")[1]
}

func getNameWithPackage(namespacedName string, level int) string {
	splitDash := strings.Split(namespacedName, "-")

	packagePath := splitDash[0]
	realName := splitDash[1]

	splitSlashes := strings.Split(packagePath, "/")

	var index = len(splitSlashes) - 1
	output := realName

	for i := index; i >= index-level+1; i-- {
		output = splitSlashes[i] + output
	}

	return output
}

func structsToValibot(structList StructList) (string, error) {
	valibotOutput := ""

	names := make([]string, len(structList))
	nameToIndex := make(map[string]int)
	usedNames := make([]string, 0)
	nameMap := make(map[string]string)

	for i, v := range structList {
		names[i] = v.Name
		nameToIndex[v.Name] = i
	}

	sort.Slice(names, func(i, j int) bool {
		n1 := names[i]
		n2 := names[j]

		return strings.Count(n1, "-")+strings.Count(n1, "/") < strings.Count(n2, "-")+strings.Count(n2, "/")
	})

	for _, name := range names {
		structName := getName(name)

		originalIndex := nameToIndex[name]
		originalStruct := structList[originalIndex]

		exists := slices.Contains(usedNames, structName)

		level := 1
		for exists {
			structName = getNameWithPackage(originalStruct.Name, level)

			level++
			exists = slices.Contains(usedNames, structName)
		}

		usedNames = append(usedNames, structName)
		nameMap[originalStruct.Name] = structName
	}

	importedValidators := make(map[string]uint)
	importedValidators["object"] = 0
	var counter uint = 1

	for _, s := range structList {
		localValidbotOutput := "const " + nameMap[s.Name] + " = object({\n"

		for _, fieldType := range s.Fields {
			fieldOutput, err := getSingleField(importedValidators, nameMap, &counter, fieldType, 0)
			if err != nil {
				return "", err
			}

			localValidbotOutput += fieldOutput
		}

		localValidbotOutput += "});"
		valibotOutput += "\n" + localValidbotOutput + "\n"
	}

	validatorsArr := make([]string, len(importedValidators))
	for k, v := range importedValidators {
		validatorsArr[v] = k
	}

	importLine := ""
	for _, validator := range validatorsArr {
		importLine += validator + ", "
	}

	importLine = "import { " + importLine[:len(importLine)-2] + " } from 'valibot';\n"

	return "\n" + importLine + valibotOutput, nil
}
