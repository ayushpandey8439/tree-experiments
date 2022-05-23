// The intervals in this are generated based on a formula for hierarchy.
// When a leaf is added, the path from the leaf to the root needs to be updated.
package main

import (
	"fmt"
	"io"
	"os"
)

type BinaryNode struct {
	left     *BinaryNode
	right    *BinaryNode
	parent   *BinaryNode
	data     int
	interval [2]int
}

type BinaryTree struct {
	root *BinaryNode
}

func (t *BinaryTree) insert(data int) *BinaryTree {
	if t.root == nil {
		t.root = &BinaryNode{data: data, left: nil, right: nil, parent: nil, interval: [2]int{2, 2}}
	} else {
		t.root.insert(data)
	}
	return t
}

func (n *BinaryNode) insert(data int) {
	if n == nil {
		return
	} else if data <= n.data {
		if n.left == nil {
			var iInter = ((n.interval[0] - 1) * 10) + 9
			n.left = &BinaryNode{data: data, left: nil, right: nil, parent: n, interval: [2]int{iInter, iInter}}
		} else {
			n.left.insert(data)
		}
	} else {
		if n.right == nil {
			var iInter = (n.interval[1] * 10) + 1
			n.right = &BinaryNode{data: data, left: nil, right: nil, parent: n, interval: [2]int{iInter, iInter}}
		} else {
			n.right.insert(data)
		}
	}

	if n.left != nil && n.right != nil {
		n.interval[0] = n.left.interval[0]
		n.interval[1] = n.right.interval[1]
	} else {
		if n.right == nil {
			n.interval = n.left.interval
		}
		if n.left == nil {
			n.interval = n.right.interval

		}
	}

}
func print(w io.Writer, node *BinaryNode, ns int, ch rune) {
	if node == nil {
		return
	}

	for i := 0; i < ns; i++ {
		fmt.Fprint(w, " ")
	}
	fmt.Fprintf(w, "%c:%v       %v\n", ch, node.data, node.interval)
	print(w, node.left, ns+2, 'L')
	print(w, node.right, ns+2, 'R')
}

func subsumes(node1 *BinaryNode, node2 *BinaryNode) {

}

func main() {
	tree := &BinaryTree{}
	tree.insert(100).
		insert(-20).
		insert(150)

	print(os.Stdout, tree.root, 0, 'M')
}
