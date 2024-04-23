package main

import (
	"errors"
	"slices"
	"sort"
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

// ==================================================
// Backend Methods.
//
// Responsible for ordering output, and producing,
// the resultant code.
// ==================================================

func recGetLastType(field FieldInfo) string {
	if field.Value == nil {
		return field.Type
	}

	return recGetLastType(*field.Value)
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

func getSingleField(validators map[string]uint, counter *uint, field FieldInfo) (string, error) {
	jsType, err := getJsType(recGetLastType(field))
	localValibotOutput := ""

	if err == nil {
		maybeAdd(validators, counter, jsType)
	}

	if err == NoJsType && !(field.Map || field.Array) {
		localValibotOutput += "  " + field.Name + ": " + field.Type + ",\n"
		return localValibotOutput, nil
	}

	localValibotOutput += "  " + field.Name + ": "

	if field.Array {
		maybeAdd(validators, counter, "array")
		localValibotOutput += "array(" + appendedType(field.Type)
		return localValibotOutput, nil
	}

	if field.Map {
		maybeAdd(validators, counter, "record")

		localValibotOutput += "record("

		// TODO: Refactor me! I am here for 1 test to pass!
		if field.Value != nil && field.Value.Array {
			maybeAdd(validators, counter, "array")
			hackType := appendedType(field.Value.Type)
			localValibotOutput += "array(" + hackType[0:len(hackType)-2] + "),\n"
		} else {
			localValibotOutput += appendedType(field.Type)
		}

		return localValibotOutput, nil
	}

	localValibotOutput += jsType + "(),\n"
	return localValibotOutput, nil
}

func structsToValibot(structList StructList) (string, error) {
	valibotOutput := ""

	importedValidators := make(map[string]uint)
	importedValidators["object"] = 0
	var counter uint = 1

	for _, s := range structList {
		localValidbotOutput := "const " + s.Name + " = object({\n"

		for _, fieldType := range s.Fields {
			fieldOutput, err := getSingleField(importedValidators, &counter, fieldType)
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
