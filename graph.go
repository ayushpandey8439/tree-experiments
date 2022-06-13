package main

// TODO: Paths of equal lengths cause problems and the result is based on the update order. This needs to be mitigated.
import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
)

var loggingMode string

type vertex struct {
	ID       int
	children map[int]*vertex
	parents  map[int]*vertex
	path     [2][]int
	data     int
}

// We only allow a single vertex entry to the graph.
type graph struct {
	nodeCounter int
	vertices    map[int]*vertex
	root        *vertex
	paths       [][]int
	currentPath []int
}

func (G *graph) createVertex(V int) *graph {
	Node := &vertex{
		ID:       G.nodeCounter,
		children: make(map[int]*vertex),
		parents:  make(map[int]*vertex),
		path:     [2][]int{{G.nodeCounter}, {G.nodeCounter}},
		data:     V,
	}
	G.vertices[G.nodeCounter] = Node
	G.nodeCounter++
	return G
}

func (G *graph) createEdge(V1 int, V2 int) *graph {
	source := G.vertices[V1]
	target := G.vertices[V2]

	source.children[V2] = target
	target.parents[V1] = source

	G.updatePath(source, target, false, [2][]int{}, source.ID)

	return G
}

func (G *graph) updatePath(source *vertex, target *vertex, isInherited bool, inheritedPath [2][]int, inheritedFrom int) *graph {

	if isInherited {
		targetPLow := make([]int, len(target.path[0]))
		copy(targetPLow, target.path[0])
		targetPHigh := make([]int, len(target.path[1]))
		copy(targetPHigh, target.path[1])

		if targetPLow[0] == inheritedFrom {
			newLow := append(inheritedPath[0][:len(inheritedPath[0])-1], targetPLow...)
			//fmt.Fprintf(os.Stdout, "\nUpdating Low Inherited Paths on %v, %v with %v", target.ID, inheritedPath, newLow)
			target.path[0] = make([]int, len(newLow))
			copy(target.path[0], newLow)
		}
		if targetPHigh[0] == inheritedFrom {
			newHigh := append(inheritedPath[1][:len(inheritedPath[1])-1], targetPHigh...)
			//fmt.Fprintf(os.Stdout, "\nUpdating High Inherited Paths on %v, %v with %v", target.ID, inheritedPath, newHigh)
			target.path[1] = make([]int, len(newHigh))
			copy(target.path[1], newHigh)
		}

	} else {
		sourcePLow := make([]int, len(source.path[0]))
		copy(sourcePLow, source.path[0])
		sourcePHigh := make([]int, len(source.path[1]))
		copy(sourcePHigh, source.path[1])

		targetPLow := make([]int, len(target.path[0]))
		copy(targetPLow, target.path[0])
		targetPHigh := make([]int, len(target.path[1]))
		copy(targetPHigh, target.path[1])

		targetPLowNew := append(sourcePLow, target.ID)
		targetPHighNew := append(sourcePHigh, target.ID)
		fmt.Fprintf(os.Stdout, "\nUpdate edge %v, %v Old Paths are %v %v ", source.ID, target.ID, targetPLow, targetPHigh)

		if len(targetPLowNew) == len(targetPLow) {
			if G.pathSmallerThan(targetPLowNew, targetPLow) {
				target.path[0] = targetPLowNew
			}
		} else if len(targetPLowNew) < len(targetPLow) || len(targetPLow) == 1 {
			target.path[0] = targetPLowNew
		}
		if len(targetPHighNew) == len(targetPHigh) {
			if G.pathSmallerThan(targetPHigh, targetPHighNew) {
				target.path[1] = targetPHighNew
			}
		} else if len(targetPHighNew) > len(targetPHigh) || len(targetPLow) == 1 {
			target.path[1] = targetPHighNew
		}
		fmt.Fprintf(os.Stdout, "Updated paths are %v %v\n", target.path[0], target.path[1])

	}

	for _, v := range target.children {
		if isInherited {
			G.updatePath(target, v, true, inheritedPath, inheritedFrom)
		} else {
			G.updatePath(target, v, true, target.path, target.ID)
		}
	}

	return G
}
func (G *graph) pathSmallerThan(P1 []int, P2 []int) bool {
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
		V := G.vertices[i]

		fmt.Fprintf(w, "Node:  %v %v\t Children: ", V.ID, &V)
		fmt.Fprintf(w, "%v \t", V.children)
		fmt.Fprintf(w, "\t Paths: %v\t \n", V.path)

	}
	fmt.Fprintf(w, "\n\n")
}

func (G *graph) findPathLCSA(V ...int) int {
	StudyPaths := make([][]int, 0)
	for _, v := range V {
		vert := G.vertices[v]
		PathLow := vert.path[0]
		PathHigh := vert.path[1]
		StudyPaths = append(StudyPaths, PathLow)
		StudyPaths = append(StudyPaths, PathHigh)
	}

	shortestPath := StudyPaths[0]
	for _, path := range StudyPaths {
		if len(path) < len(shortestPath) {
			shortestPath = path
		}
	}

	lowestNode := 0
	for i := 0; i < len(shortestPath); i++ {
		node := shortestPath[i]
		AllHaveNode := true
		for j := 0; j < len(StudyPaths); j++ {
			if StudyPaths[j][i] != node {
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
	PLCSA := &vertex{ID: 0}
	if lowestNode == 0 {
		PLCSA = nil
	} else {
		PLCSA = G.vertices[lowestNode]
	}

	if loggingMode != "none" {
		fmt.Fprintf(os.Stdout, "LCSA via longest path prefix is %v \n", PLCSA.ID)
	}

	return PLCSA.ID
}

func (G *graph) findTraversalLCSA(V ...int) int {
	if loggingMode != "none" {
		fmt.Fprintf(os.Stdout, "\nPaths to nodes ")
	}
	for _, v := range V {
		if loggingMode != "none" {
			fmt.Fprintf(os.Stdout, " %v", v)
		}

		G.dfs(G.root.ID, v)
	}
	if loggingMode != "none" {
		fmt.Fprintf(os.Stdout, "\n")
	}

	if loggingMode != "none" {
		for i := 0; i < len(G.paths); i++ {
			fmt.Fprintf(os.Stdout, "\t%s \n", G.paths[i])
		}
	}

	shortestPath := []int{}
	for i := 0; i < len(G.paths); i++ {
		if len(G.paths[i]) <= len(shortestPath) || len(shortestPath) == 0 {
			shortestPath = G.paths[i]
		}
	}
	lowestNode := 0
	for i := 0; i < len(shortestPath); i++ {
		node := shortestPath[i]
		AllHaveNode := true
		for j := 0; j < len(G.paths); j++ {
			if G.paths[j][i] != node {
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
	G.paths = make([][]int, 0)
	return G.vertices[lowestNode].ID

}

func (G *graph) dfs(source int, dest int) {
	if source == dest {
		detectedpath := make([]int, len(G.currentPath))
		copy(detectedpath, G.currentPath)
		G.paths = append(G.paths, detectedpath)
		return
	} else {
		for _, V := range G.vertices[source].children {
			G.currentPath = append(G.currentPath, V.ID)
			G.dfs(V.ID, dest)
			G.currentPath = G.currentPath[:len(G.currentPath)-1]
		}
	}
}

func main() {
	loggingMode = "none"
	G := &graph{
		vertices:    make(map[int]*vertex),
		nodeCounter: 1,
	}

	// Define the number of nodes here.
	numNodes := 10
	for j := 1; j <= numNodes; j++ {
		G.createVertex(j)
	}

	G.root = G.vertices[1]
	G.currentPath = []int{G.root.ID}

	edgeMap := make(map[int][]int)
	edgeMap[1] = []int{2, 3}
	edgeMap[2] = []int{4, 5}
	edgeMap[3] = []int{4, 5, 6, 7}
	edgeMap[6] = []int{9}
	edgeMap[7] = []int{8, 9, 10}
	edgeMap[8] = []int{10}

	for Source, Targets := range edgeMap {
		for _, Target := range Targets {
			//fmt.Fprintf(os.Stdout, "\nCreating edge between %v, %v \n", Source, Target)
			G.createEdge(Source, Target)
		}
	}

	printGraph(os.Stdout, G)
	fmt.Fprintf(os.Stdout, "\nAll Pair LCSA failed for  %v \n\n", G.testAllPairLCSA(numNodes))
	runtime.GC()
}

func (G *graph) testLCSA(v ...int) {
	fmt.Fprintf(os.Stdout, "\nLCSA same via path and dfs?  %v \n\n", G.findPathLCSA(v...) == G.findTraversalLCSA(v...))
}

func (G *graph) testAllPairLCSA(n int) []string {
	FailedPairs := []string{}
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			//G.findPathLCSA(i, j)
			PLCSA := G.findPathLCSA(i, j)
			TLCSA := G.findTraversalLCSA(i, j)
			if !(PLCSA == TLCSA) {
				FailedPairs = append(FailedPairs, strconv.Itoa(i)+","+strconv.Itoa(j)+"; PLCSA: "+strconv.Itoa(PLCSA)+"; TLCSA: "+strconv.Itoa(TLCSA))
			}
		}
	}
	return FailedPairs
}
