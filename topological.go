package main

import "slices"

type Node struct {
	Name    string
	Visited bool
	Edges   []Node
}

func dfs(node Node) []string {
	ordering := make([]string, 0)

	for _, n := range node.Edges {
		if n.Visited {
			continue
		}

		ordering = slices.Concat(ordering, dfs(n))
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
