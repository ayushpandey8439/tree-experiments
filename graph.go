package main

// TODO: Paths of equal lengths cause problems and the result is based on the update order. This needs to be mitigated.
import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var loggingMode string

type vertex struct {
	ID       string
	children map[string]*vertex
	parents  map[string]*vertex
	path     [2]string
	data     int
}

// We only allow a single vertex entry to the graph.
type graph struct {
	nodeCounter int
	vertices    map[string]*vertex
	root        *vertex
	paths       [][]string
}

func length(path string) int {
	return len(strings.Split(path, ";"))
}

func (G *graph) createVertex(V int) *graph {
	var nodeCounterString = strconv.Itoa(G.nodeCounter)
	Node := &vertex{
		ID:       nodeCounterString,
		children: make(map[string]*vertex),
		parents:  make(map[string]*vertex),
		path:     [2]string{nodeCounterString, nodeCounterString},
		data:     V,
	}
	G.vertices[nodeCounterString] = Node
	G.nodeCounter++
	return G
}

func (G *graph) createEdge(V1 int, V2 int) *graph {
	v1 := strconv.Itoa(V1)
	v2 := strconv.Itoa(V2)

	source := G.vertices[v1]
	target := G.vertices[v2]

	source.children[v2] = target
	target.parents[v1] = source

	G.updatePath(source, target, false)

	return G
}

func (G *graph) removeEdge(V1 int, V2 int) *graph {
	v1 := strconv.Itoa(V1)
	v2 := strconv.Itoa(V2)

	source := G.vertices[v1]
	target := G.vertices[v2]

	delete(source.children, v2)
	delete(target.parents, v1)

	for _, newParent := range target.parents {
		if len(newParent.path[0]) < len(target.path[0]) {
			target.path[0] = newParent.path[0]
		}
		if len(newParent.path[1]) > len(target.path[1]) {
			target.path[1] = newParent.path[1]
		}
	}

	for _, v := range target.children {
		G.updatePath(target, v, true) // the inherited path is used to update all the paths of a subtree
	}
	return G
}

func (G *graph) updatePath(source *vertex, target *vertex, isInherited bool) *graph {
	sourcePLow := source.path[0]
	sourcePHigh := source.path[1]

	targetPLow := target.path[0]
	targetPHigh := target.path[1]

	targetPLowNew := sourcePLow + ";" + target.ID
	targetPHighNew := sourcePHigh + ";" + target.ID

	if len(targetPLowNew) == len(targetPLow) || isInherited {
		fmt.Fprintf(os.Stdout, "\nPath Length Conflict for %v \n", target.ID)
		if G.pathSmallerThan(targetPLowNew, targetPLow) {
			fmt.Fprintf(os.Stdout, "\nBroken conflict and shorter path is %v \n", targetPLowNew)
			target.path[0] = targetPLowNew
		}
	} else if len(targetPLowNew) < len(targetPLow) || targetPLow == target.ID || isInherited {
		fmt.Fprintf(os.Stdout, "\nFound a shorter path to %v \n", target.ID)
		target.path[0] = targetPLowNew
	}
	if len(targetPHighNew) == len(targetPHigh) {
		fmt.Fprintf(os.Stdout, "\nPath Length Conflict for %v \n", target.ID)
		if G.pathSmallerThan(targetPHigh, targetPHighNew) {
			fmt.Fprintf(os.Stdout, "\nBroken conflict and longer path is %v \n", targetPHighNew)
			target.path[1] = targetPHighNew
		}
	} else if len(targetPHighNew) > len(targetPHigh) || targetPLow == target.ID || isInherited {
		fmt.Fprintf(os.Stdout, "\nFound a longer path to %v \n", target.ID)
		target.path[1] = targetPHighNew
	}

	for _, v := range target.children {
		G.updatePath(target, v, true)
		// The inherited path is used to update all the paths of a subtree.
		// This inherited property comes of importance when the root of the tree changes and
		// a new root is added at a level above the current root
	}

	return G
}
func (G *graph) pathSmallerThan(P1 string, P2 string) bool {
	// This method assumes that for any node, the children are ordered by their Ids. So a child with a highest ID is the right most child
	smallerPath := false
	for i := 0; i < len(P1); i++ {
		node := P1[i]
		if P2[i] != node {
			lowestNode1 := P1[i]
			lowestNode2 := P2[i]
			if lowestNode1 < lowestNode2 {
				smallerPath = true
			}
		}
	}
	return smallerPath
}

func printGraph(w io.Writer, G *graph) {

	for i := 1; i <= len(G.vertices); i++ {
		V := G.vertices[strconv.Itoa(i)]

		fmt.Fprintf(w, "Node:  %v\t Children: ", V.ID)
		for _, child := range V.children {
			fmt.Fprintf(w, "%v ", child.ID)
		}
		fmt.Fprintf(w, "\t Paths: %v\t \n", V.path)

	}
	fmt.Fprintf(w, "\n\n")
}

func (G *graph) findPathLCSA(V ...int) string {
	StudyPaths := [][]string{}
	for _, v := range V {
		vs := strconv.Itoa(v)
		vert := G.vertices[vs]
		PathLow := strings.Split(vert.path[0], ";")
		PathHigh := strings.Split(vert.path[1], ";")
		StudyPaths = append(StudyPaths, PathLow)
		StudyPaths = append(StudyPaths, PathHigh)
	}

	shortestPath := StudyPaths[0]

	for _, path := range StudyPaths {
		if len(path) < len(shortestPath) {
			shortestPath = path
		}
	}

	lowestNode := ""
	for i := 0; i < len(shortestPath); i++ {
		node := string(shortestPath[i])
		AllHaveNode := true
		for j := 0; j < len(StudyPaths); j++ {
			if string(StudyPaths[j][i]) != node {
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

	PLCSA := G.vertices[lowestNode]

	if loggingMode != "none" {
		fmt.Fprintf(os.Stdout, "LCSA via longest path prefix is %v \n", PLCSA.ID)
	}

	return PLCSA.ID
}

func (G *graph) findTraversalLCSA(V ...int) string {
	if loggingMode != "none" {
		fmt.Fprintf(os.Stdout, "\nPaths to nodes ")
	}
	for _, v := range V {
		vs := strconv.Itoa(v)
		if loggingMode != "none" {
			fmt.Fprintf(os.Stdout, " %v", vs)
		}

		G.dfs(G.root.ID, vs, make(map[string]bool), G.root.ID)
	}
	if loggingMode != "none" {
		fmt.Fprintf(os.Stdout, "\n")
	}

	if loggingMode != "none" {
		for i := 0; i < len(G.paths); i++ {
			fmt.Fprintf(os.Stdout, "\t%s \n", G.paths[i])
		}
	}

	shortestPath := []string{}
	for i := 0; i < len(G.paths); i++ {
		if len(G.paths[i]) <= len(shortestPath) || len(shortestPath) == 0 {
			shortestPath = G.paths[i]
		}
	}
	lowestNode := ""
	for i := 0; i < len(shortestPath); i++ {
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

	if loggingMode != "none" {
		fmt.Fprintf(os.Stdout, "\nLCSA via Path Traversal is %v \n", lowestNode)
	}
	G.paths = [][]string{}
	return G.vertices[lowestNode].ID

}

func (G *graph) dfs(source string, dest string, visiting map[string]bool, currentPath string) {
	if source == dest {
		currPath := strings.Split(currentPath, ";")
		G.paths = append(G.paths, currPath)
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

func main() {
	loggingMode = "all"
	G := &graph{
		vertices:    make(map[string]*vertex),
		nodeCounter: 1,
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
		for _, Target := range Targets {
			fmt.Fprintf(os.Stdout, "\nCreating edge between %v, %v \n", Source, Target)
			G.createEdge(Source, Target)
		}
	}

	printGraph(os.Stdout, G)
	G.testLCSA(7, 9)
	//G.testAllPairLCSA(numNodes)
	//G.removeEdge(7, 10)
	//G.testAllPairLCSA(numNodes)

	G = &graph{
		vertices:    make(map[string]*vertex),
		nodeCounter: 1,
	}
}

func (G *graph) testLCSA(v ...int) {
	fmt.Fprintf(os.Stdout, "\nLCSA same via path and dfs?  %v \n\n", G.findPathLCSA(v...) == G.findTraversalLCSA(v...))
}

func (G *graph) testAllPairLCSA(n int) {
	FailedPairs := []string{}
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if !(G.findPathLCSA(i, j) == G.findTraversalLCSA(i, j)) {
				FailedPairs = append(FailedPairs, strconv.Itoa(i)+","+strconv.Itoa(j))
			}
		}
	}
	fmt.Fprintf(os.Stdout, "\nAll Pair LCSA failed for  %v \n\n", FailedPairs)
}
