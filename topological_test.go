package main

import (
	"slices"
	"testing"
)

func TestSingle(t *testing.T) {
	t.Skip()
	node := Node{
		Name:    "A",
		Visited: false,
		Edges:   make([]*Node, 0),
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
	t.Skip()
	n1 := Node{
		Name:    "A",
		Visited: false,
		Edges:   make([]*Node, 0),
	}

	n2 := Node{
		Name:    "B",
		Visited: false,
		Edges:   make([]*Node, 0),
	}

	n1.Edges = append(n1.Edges, &n2)

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

func TestComplex(t *testing.T) {
	n1 := Node{
		Name:    "A",
		Visited: false,
		Edges:   make([]*Node, 0),
	}

	n2 := Node{
		Name:    "B",
		Visited: false,
		Edges:   make([]*Node, 0),
	}

	n3 := Node{
		Name:    "C",
		Visited: false,
		Edges:   make([]*Node, 0),
	}

	n4 := Node{
		Name:    "D",
		Visited: false,
		Edges:   make([]*Node, 0),
	}

	n5 := Node{
		Name:    "F",
		Visited: false,
		Edges:   make([]*Node, 0),
	}

	n1.Edges = append(n1.Edges, &n2, &n3)
	n2.Edges = append(n2.Edges, &n3, &n4)
	n3.Edges = append(n3.Edges, &n4, &n5)
	n4.Edges = append(n4.Edges, &n5)

	correctOrder := []string{"F", "D", "C", "B", "A"}

	nodeSlice := []Node{n3, n4, n2, n1, n5}
	testOrder := topologicalSort(nodeSlice)

	t.Log(testOrder)

	if slices.Compare(correctOrder, testOrder) != 0 {
		t.Log("Slices are not equal")
		t.FailNow()
	}
}
