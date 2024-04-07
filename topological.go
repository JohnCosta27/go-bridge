package main

type Node struct {
	Name    string
	Visited bool
	Edges   []*Node
}

func dfs(node *Node, list *[]string) *[]string {
	if node.Visited {
		return list
	}

	for _, n := range node.Edges {
		dfs(n, list)
	}

	node.Visited = true
	*list = append(*list, node.Name)

	return list
}

func hasUnvisitedNode(nodes []*Node) bool {
	for _, n := range nodes {
		if !n.Visited {
			return true
		}
	}

	return false
}

func topologicalSort(nodes []*Node) []string {
	longest := []string{}

	for hasUnvisitedNode(nodes) {
		for _, n := range nodes {
			if !n.Visited {
				dfs(n, &longest)
			}
		}
	}

	return longest
}
