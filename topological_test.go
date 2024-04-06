package main

import (
	"slices"
	"testing"
)

func TestSingle(t *testing.T) {
	node := Node{
		Name:    "A",
		Visited: false,
		Edges:   make([]Node, 0),
	}

	correctOrder := []string{"A"}

	nodeSlice := []Node{node}
	testOrder := topologicalSort(nodeSlice)

	t.Log(testOrder)

	if slices.Compare(correctOrder, testOrder) != 0 {
		t.Log("Slices are not equal")
		t.FailNow()
	}
}

func TestTwo(t *testing.T) {
	n1 := Node{
		Name:    "A",
		Visited: false,
		Edges:   make([]Node, 0),
	}

	n2 := Node{
		Name:    "B",
		Visited: false,
		Edges:   make([]Node, 0),
	}

	n1.Edges = append(n1.Edges, n2)

	correctOrder := []string{"B", "A"}

	nodeSlice := []Node{n1, n2}
	testOrder := topologicalSort(nodeSlice)

	t.Log(testOrder)

	if slices.Compare(correctOrder, testOrder) != 0 {
		t.Log("Slices are not equal")
		t.FailNow()
	}

	nodeSlice = []Node{n2, n1}
	testOrder = topologicalSort(nodeSlice)

	t.Log(testOrder)

	if slices.Compare(correctOrder, testOrder) != 0 {
		t.Log("Slices are not equal")
		t.FailNow()
	}
}
