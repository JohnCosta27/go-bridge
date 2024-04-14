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
			_, err := getJsType(field.Type)

			if err != NoJsType {
				continue
			}

			nodeIndex := slices.IndexFunc(nodeMap, func(n *Node) bool {
				return n.Name == field.Type
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

func structsToValibot(structList StructList) (string, error) {
	valibotOutput := ""

	importedValidators := make([]string, 0)
	importedValidators = append(importedValidators, "object")

	for _, s := range structList {
		localValidbotOutput := "const " + s.Name + " = object({\n"

		for _, fieldType := range s.Fields {
			jsType, err := getJsType(fieldType.Type)

			if err == NoJsType {
				// Here, we must have a nested struct.
				localValidbotOutput += "  " + fieldType.Name + ": " + fieldType.Type + ",\n"

				continue
			}

			exist := slices.Index(importedValidators, jsType) != -1
			if !exist {
				importedValidators = append(importedValidators, jsType)
			}
			localValidbotOutput += "  " + fieldType.Name + ": " + jsType + "(),\n"

		}

		localValidbotOutput += "});"
		valibotOutput += "\n" + localValidbotOutput + "\n"
	}

	importLine := ""
	for _, validator := range importedValidators {
		importLine += validator + ", "
	}

	importLine = "import { " + importLine[:len(importLine)-2] + " } from 'valibot';\n"

	return "\n" + importLine + valibotOutput, nil
}
