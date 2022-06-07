package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var nodeCounter int

type vertex struct {
	ID       string
	children []*vertex
	parents  []*vertex
	path     [2]string
	data     int
}

// We only allow a single vertex entry to the graph.
type graph struct {
	vertices map[string]*vertex
	root     *vertex
	paths    []string
}

func (G *graph) updatePath(source *vertex, target *vertex, isInherited bool) *graph {

	sourcePLow := source.path[0]
	sourcePHigh := source.path[1]

	targetPLow := target.path[0]
	targetPHigh := target.path[1]

	targetPLowNew := sourcePLow + ";" + target.ID
	targetPHighNew := sourcePHigh + ";" + target.ID

	if len(targetPLowNew) <= length(targetPLow) || targetPLow == target.ID || isInherited {
		target.path[0] = targetPLowNew
	}
	if len(targetPHighNew) >= length(targetPHigh) || targetPLow == target.ID || isInherited {
		target.path[1] = targetPHighNew
	}

	for i := 0; i < len(target.children); i++ {
		G.updatePath(target, target.children[i], true) // the inherited path is used to update all the paths of a subtree
	}

	return G
}

func length(path string) int {
	return len(strings.Split(path, ";"))
}

func (G *graph) createVertex(V int) *graph {
	var nodeCounterString = strconv.Itoa(nodeCounter)
	Node := &vertex{
		ID:       nodeCounterString,
		children: nil,
		parents:  nil,
		path:     [2]string{nodeCounterString, nodeCounterString},
		data:     V,
	}
	G.vertices[nodeCounterString] = Node
	nodeCounter++
	return G
}

func (G *graph) createEdge(V1 int, V2 int) *graph {
	v1 := strconv.Itoa(V1)
	v2 := strconv.Itoa(V2)

	source := G.vertices[v1]
	target := G.vertices[v2]

	source.children = append(source.children, target)
	target.parents = append(target.parents, source)

	G.updatePath(source, target, false)

	return G
}

func (G *graph) findPathLCSA(V1 int, V2 int) {
	v1 := strconv.Itoa(V1)
	v2 := strconv.Itoa(V2)

	A := G.vertices[v1]
	B := G.vertices[v2]

	APathLow := strings.Split(A.path[0], ";")
	APathHigh := strings.Split(A.path[1], ";")

	shorterAPath := APathHigh
	if len(APathLow) < len(APathHigh) {
		shorterAPath = APathLow
	}

	ACommonPath := ""

	for i := 0; i < len(shorterAPath); i++ {
		if APathLow[i] == APathHigh[i] {
			ACommonPath = string(ACommonPath) + string(APathLow[i])
		} else {
			break
		}
	}

	BPathLow := strings.Split(B.path[0], ";")
	BPathHigh := strings.Split(B.path[1], ";")

	shorterBPath := BPathHigh
	if len(BPathLow) < len(BPathHigh) {
		shorterBPath = BPathLow
	}

	BCommonPath := ""

	for i := 0; i < len(shorterBPath); i++ {
		if BPathLow[i] == BPathHigh[i] {
			BCommonPath = string(BCommonPath) + string(BPathLow[i])
		} else {
			break
		}
	}

	PLCSA := &vertex{}
	shorterPath := BCommonPath
	if len(ACommonPath) < len(BCommonPath) {
		shorterPath = ACommonPath
	}

	for i := 0; i < len(shorterPath); i++ {
		if ACommonPath[i] == BCommonPath[i] {
			PLCSA = G.vertices[string(ACommonPath[i])]
		} else {
			break
		}
	}

	fmt.Fprintf(os.Stdout, "LCSA via longest path prefix is %v \n\n", PLCSA.ID)
}

func (G *graph) findTraversalLCSA(V1 int, V2 int) {
	v1 := strconv.Itoa(V1)
	v2 := strconv.Itoa(V2)

	G.dfs(G.root.ID, v1, make(map[string]bool), G.root.ID)

	G.dfs(G.root.ID, v2, make(map[string]bool), G.root.ID)

	fmt.Fprintf(os.Stdout, "Paths to %v or %v are:: \n", v1, v2)
	for i := 0; i < len(G.paths); i++ {
		fmt.Fprintf(os.Stdout, "\t%s \n", G.paths[i])
	}

	shortestPath := ""
	for i := 0; i < len(G.paths); i++ {
		if len(G.paths[i]) > len(shortestPath) || len(shortestPath) == 0 {
			shortestPath = G.paths[i]
		}
	}
	lowestNode := ""
	for i := 0; i < len(shortestPath); i = i + 2 {
		node := string(shortestPath[i])
		AllHaveNode := true
		for j := 0; j < len(G.paths); j++ {
			if string(G.paths[j][i]) != node {
				AllHaveNode = false
				break
			}
		}
		if !AllHaveNode {
			break
		} else {
			lowestNode = node
		}
	}

	fmt.Fprintf(os.Stdout, "\nLCSA via Path Traversal is %v \n\n", lowestNode)

}

func (G *graph) dfs(source string, dest string, visiting map[string]bool, currentPath string) {
	if source == dest {
		G.paths = append(G.paths, currentPath)
		return
	}

	visiting[source] = true

	for _, V := range G.vertices[source].children {
		if !visiting[V.ID] {
			currentPath = currentPath + ";" + V.ID
			G.dfs(V.ID, dest, visiting, currentPath)
			currentPath = strings.TrimSuffix(currentPath, ";"+V.ID)
		}
	}
}

func printGraph(w io.Writer, G *graph) {

	for i := 1; i <= len(G.vertices); i++ {
		V := G.vertices[strconv.Itoa(i)]

		fmt.Fprintf(w, "Node:  %v\t Children: ", V.ID)
		for j := 0; j < len(V.children); j++ {
			fmt.Fprintf(w, "%v ", V.children[j].ID)
		}
		fmt.Fprintf(w, "\t Paths: %v\t \n", V.path)

	}
	fmt.Fprintf(w, "\n\n")
}

func main() {
	nodeCounter = 1

	G := &graph{
		vertices: make(map[string]*vertex),
	}

	// Define the number of nodes here.
	numNodes := 10
	for j := 1; j <= numNodes; j++ {
		G.createVertex(j)
	}

	G.root = G.vertices["1"]

	edgeMap := make(map[int][]int)
	edgeMap[1] = []int{2, 3}
	edgeMap[2] = []int{4, 5}
	edgeMap[3] = []int{4, 5, 6, 7}
	edgeMap[6] = []int{9}
	edgeMap[7] = []int{8, 9, 10}
	edgeMap[8] = []int{10}

	for Source, Targets := range edgeMap {
		for i := 0; i < len(Targets); i++ {
			G.createEdge(Source, Targets[i])
		}
	}
	printGraph(os.Stdout, G)

	G.findPathLCSA(5, 7)
	G.findTraversalLCSA(5, 7)
}
