// The intervals in this are generated based on a formula for hierarchy.
// When a leaf is added, the path from the leaf to the root needs to be updated.
package main

import (
	"fmt"
	"io"
	"os"
)

type PositionBinaryNode struct {
	left     *PositionBinaryNode
	right    *PositionBinaryNode
	parent   *PositionBinaryNode
	data     int
	position string
}

type PositionBinaryTree struct {
	root *PositionBinaryNode
}

func (t *PositionBinaryTree) insert(data int) *PositionBinaryTree {
	if t.root == nil {
		t.root = &PositionBinaryNode{data: data, left: nil, right: nil, parent: nil, position: "t"}
	} else {
		t.root.insert(data)
	}
	return t
}

func (n *PositionBinaryNode) insert(data int) {
	if n == nil {
		return
	} else if data <= n.data {
		if n.left == nil {
			pos := n.position + "l"
			n.left = &PositionBinaryNode{data: data, left: nil, right: nil, parent: n, position: pos}
		} else {
			n.left.insert(data)
		}
	} else {
		if n.right == nil {
			pos := n.position + "r"
			n.right = &PositionBinaryNode{data: data, left: nil, right: nil, parent: n, position: pos}
		} else {
			n.right.insert(data)
		}
	}
}

func positionprint(w io.Writer, node *PositionBinaryNode, ns int, ch rune) {
	if node == nil {
		return
	}

	for i := 0; i < ns; i++ {
		fmt.Fprint(w, " ")
	}
	fmt.Fprintf(w, "%c:%v       %v\n", ch, node.data, node.position)
	positionprint(w, node.left, ns+2, 'L')
	positionprint(w, node.right, ns+2, 'R')
}

func main() {
	tree := &PositionBinaryTree{}

	tree.insert(100).
		insert(-20).
		insert(150)

	positionprint(os.Stdout, tree.root, 0, 'M')
}
