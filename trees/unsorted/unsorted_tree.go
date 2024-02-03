package unsorted_tree

import (
	"fmt"
	"strings"
)

type Node struct {
	data        string
	left, right *Node
}

// *** DoublyLinkedList code ***
type Cell struct {
	data       *Node
	next, prev *Cell
}

type DoublyLinkedList struct {
	topSentinel    *Cell
	bottomSentinel *Cell
}

func makeDoublyLinkedList() DoublyLinkedList {
	list := DoublyLinkedList{}
	list.topSentinel = &Cell{&Node{data: "TOP_SENTINEL"}, nil, nil}
	list.bottomSentinel = &Cell{&Node{data: "BOTTOM_SENTINEL"}, nil, nil}
	list.topSentinel.next = list.bottomSentinel
	list.bottomSentinel.prev = list.topSentinel
	return list
}

func (list *DoublyLinkedList) isEmpty() bool {
	return list.topSentinel.next == list.bottomSentinel
}

// Add a cell after me.
func (me *Cell) addAfter(after *Cell) {
	after.next = me.next
	after.prev = me
	if me.next != nil {
		me.next.prev = after
	}
	me.next = after
}

// Add an item to the top of the queue.
func (queue *DoublyLinkedList) enqueue(node *Node) {
	queue.topSentinel.addAfter(&Cell{data: node, next: nil, prev: nil})
}

// Remove an item from the bottom of the queue.
func (queue *DoublyLinkedList) dequeue() *Node {
	r := &Node{}
	if queue.bottomSentinel.prev != nil {
		r = queue.bottomSentinel.prev.data
	}
	queue.bottomSentinel.deleteBefore()
	return r
}

func (me *Cell) deleteBefore() {
	if me.prev != nil {
		me.prev.delete()
	}
}

// Delete a cell after me.
func (me *Cell) delete() {
	if me.data.data == "TOP_SENTINEL" || me.data.data == "BOTTOM_SENTINEL" {
		panic("Can't delete a sentinel")
	}
	me.prev.next = me.next
	me.next.prev = me.prev
}

func buildTree() *Node {
	aNode := Node{"A", nil, nil}
	bNode := Node{"B", nil, nil}
	cNode := Node{"C", nil, nil}
	dNode := Node{"D", nil, nil}
	eNode := Node{"E", nil, nil}
	fNode := Node{"F", nil, nil}
	gNode := Node{"G", nil, nil}
	hNode := Node{"H", nil, nil}
	iNode := Node{"I", nil, nil}
	jNode := Node{"J", nil, nil}

	aNode.left = &bNode
	aNode.right = &cNode
	bNode.left = &dNode
	bNode.right = &eNode
	eNode.left = &gNode
	cNode.right = &fNode
	fNode.left = &hNode
	hNode.left = &iNode
	hNode.right = &jNode

	return &aNode
}

func (node *Node) displayIndented(indent string, depth int) string {
	result := strings.Repeat(indent, depth) + node.data + "\n"
	if node.left != nil {
		result += node.left.displayIndented(indent, depth+1)
	}
	if node.right != nil {
		result += node.right.displayIndented(indent, depth+1)
	}
	return result
}

func (node *Node) preorder() string {
	result := node.data
	if node.left != nil {
		result += " " + node.left.preorder()
	}
	if node.right != nil {
		result += " " + node.right.preorder()
	}
	return result
}

func (node *Node) inorder() string {
	result := ""
	if node.left != nil {
		result += node.left.inorder() + " "
	}
	result += node.data
	if node.right != nil {
		result += " " + node.right.inorder()
	}
	return result
}

func (node *Node) postorder() string {
	result := ""
	if node.left != nil {
		result += node.left.postorder() + " "
	}
	if node.right != nil {
		result += node.right.postorder() + " "
	}
	result += node.data
	return result
}

func (root *Node) breadthFirst() string {
	result := ""

	// Make a queue and add the root node.
	queue := makeDoublyLinkedList()
	queue.enqueue(root)

	for !queue.isEmpty() {
		node := queue.dequeue()
		result += node.data
		if node.left != nil {
			queue.enqueue(node.left)
		}
		if node.right != nil {
			queue.enqueue(node.right)
		}
		if !queue.isEmpty() {
			result += " "
		}
	}

	return result
}

func UnsortedTreeRun() {
	// Build a tree.
	aNode := buildTree()

	// Display with indentation.
	fmt.Println(aNode.displayIndented("  ", 0))

	// Display traversals.
	fmt.Println("Preorder:     ", aNode.preorder())
	fmt.Println("Inorder:      ", aNode.inorder())
	fmt.Println("Postorder:    ", aNode.postorder())
	fmt.Println("Breadth first:", aNode.breadthFirst())
}
