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

	fmt.Fprintf(os.Stdout, "LCSA via longest path prefix is %v", PLCSA.ID)
}

func (G *graph) findLCSA(V1 int, V2 int) {
	v1 := strconv.Itoa(V1)
	v2 := strconv.Itoa(V2)

	A := G.vertices[v1]
	B := G.vertices[v2]
	AAncestors := G.findAncestors(A, make(map[string]*vertex))
	BAncestors := G.findAncestors(B, make(map[string]*vertex))
	delete(AAncestors, v1)
	delete(BAncestors, v2)

	// fmt.Fprintf(os.Stdout, "\nAncestors of %s :  %v ", v1, AAncestors)
	// fmt.Fprintf(os.Stdout, "\nAncestors of %s :  %v \n", v2, BAncestors)

	CommonAncestors := intersection(AAncestors, BAncestors)
	// fmt.Fprintf(os.Stdout, "\nCommon Ancestors :  %v \n", CommonAncestors)

	TLCSA := &vertex{}
	for _, Vertex := range CommonAncestors {
		hasNoChild := true
		for i := 0; i < len(Vertex.children); i++ {
			TestChild := Vertex.children[i].ID
			if CommonAncestors[TestChild] != nil {
				hasNoChild = false
			}
		}
		if hasNoChild {
			TLCSA = Vertex
		}
	}

	fmt.Fprintf(os.Stdout, "\n\nLCSA via traversal is %v \n\n", TLCSA.ID)

}

func intersection(m1 map[string]*vertex, m2 map[string]*vertex) map[string]*vertex {
	res := make(map[string]*vertex)
	for k, v := range m1 {
		if m2[k] != nil {
			res[k] = v
		}
	}
	return res
}

func (G *graph) findAncestors(v *vertex, ancestors map[string]*vertex) map[string]*vertex {
	ancestors[v.ID] = v
	for i := 0; i < len(v.parents); i++ {
		G.findAncestors(v.parents[i], ancestors)
	}
	return ancestors
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

	edgeMap := make(map[int][]int)
	edgeMap[1] = []int{2, 3, 6}

	for Source, Targets := range edgeMap {
		for i := 0; i < len(Targets); i++ {
			G.createEdge(Source, Targets[i])
		}
	}
	printGraph(os.Stdout, G)

	G.findPathLCSA(2, 3)
	G.findLCSA(2, 3)
}

// Add tests for checking if the LCSA is actually correct
