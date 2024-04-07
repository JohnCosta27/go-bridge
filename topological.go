package main

import (
	"slices"
)

type Node struct {
	Name    string
	Visited bool
	Edges   []*Node
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func removeDuplicates(strList []string) []string {
	list := []string{}
	for _, item := range strList {
		if contains(list, item) == false {
			list = append(list, item)
		}
	}
	return list
}

func dfs(node Node) []string {
	ordering := make([]string, 0)

	for _, n := range node.Edges {
		if n.Visited {
			continue
		}

		ordering = removeDuplicates(slices.Concat(ordering, dfs(*n)))
	}

	node.Visited = true

	ordering = append(ordering, node.Name)
	return ordering
}

func topologicalSort(nodes []Node) []string {
	longest := []string{}

	for _, n := range nodes {
		order := dfs(n)
		if len(order) > len(longest) {
			longest = order
		}
	}

	return longest
}
