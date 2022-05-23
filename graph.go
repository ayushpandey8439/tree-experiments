package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

var nodeCounter int

type vertex struct {
	ID       string
	children []*vertex
	parent   []*vertex
	path     [2]string
	data     int
}

// We only allow a single vertex entry to the graph.
type graph struct {
	vertices map[string]*vertex
}

func (G *graph) updatePath(source *vertex, target *vertex) *graph {

	sourcePLow := source.path[0]
	sourcePHigh := source.path[1]

	targetPLow := target.path[0]
	targetPHigh := target.path[1]

	targetPLowNew := sourcePLow + target.ID
	targetPHighNew := sourcePHigh + target.ID

	if len(targetPLowNew) <= len(targetPLow) || targetPLow == target.ID {
		target.path[0] = targetPLowNew
	}
	if len(targetPHighNew) >= len(targetPHigh) || targetPLow == target.ID {
		target.path[1] = targetPHighNew
	}

	for i := 0; i < len(target.children); i++ {
		G.updatePath(target, target.children[i])
	}

	return G
}

func (G *graph) createVertex(V int) *graph {
	var nodeCounterString = strconv.Itoa(nodeCounter)
	Node := &vertex{
		ID:       nodeCounterString,
		children: nil,
		parent:   nil,
		path:     [2]string{nodeCounterString, nodeCounterString},
		data:     V,
	}
	G.vertices[nodeCounterString] = Node
	nodeCounter++
	return G
}

func (G *graph) createEdge(V1 string, V2 string) *graph {
	source := G.vertices[V1]
	target := G.vertices[V2]

	source.children = append(source.children, target)
	target.parent = append(target.parent, source)

	G.updatePath(source, target)

	return G
}

func printGraph(w io.Writer, G *graph) {

	for i := 1; i <= len(G.vertices); i++ {
		fmt.Fprintf(w, "Node:  %v \n", G.vertices[strconv.Itoa(i)])
	}
	fmt.Fprintf(w, "\n\n")
}

func main() {
	nodeCounter = 1

	G := &graph{
		vertices: make(map[string]*vertex),
	}

	G.createVertex(1)
	G.createVertex(2)
	G.createVertex(3)
	G.createVertex(4)
	G.createVertex(5)

	G.createEdge("1", "2")
	G.createEdge("1", "3")
	G.createEdge("1", "4")
	G.createEdge("4", "5")
	G.createEdge("2", "3")
	G.createEdge("3", "5")
	printGraph(os.Stdout, G)
}
