package sorted_tree

import "fmt"

type Node struct {
	data        string
	left, right *Node
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

func (node *Node) insertValue(newData string) {
	newNode := Node{newData, nil, nil}
	for {
		if newData <= node.data {
			// Add here or move down the left branch.
			if node.left != nil {
				node = node.left
			} else {
				node.left = &newNode
				return
			}
		} else {
			// Add here or move down the right branch.
			if node.right != nil {
				node = node.right
			} else {
				node.right = &newNode
				return
			}
		}
	}
}

func (node *Node) findValue(value string) *Node {
	for {
		if node.data == value {
			return node
		}
		if value < node.data {
			if node.left != nil {
				node = node.left
			} else {
				return nil
			}
		} else {
			if node.right != nil {
				node = node.right
			} else {
				return nil
			}
		}
	}
}

func (node *Node) isSorted() bool {
	isSorted := true
	if node.left != nil {
		if node.data > node.left.data {
			isSorted = node.left.isSorted()
		} else {
			return false
		}
	}
	if isSorted && node.right != nil {
		if node.data < node.right.data {
			isSorted = node.right.isSorted()
		}

	}
	return isSorted
}

func SortedTreeRun() {
	// Make a root node to act as sentinel.
	root := Node{"", nil, nil}

	// Add some values.
	root.insertValue("I")
	root.insertValue("G")
	root.insertValue("C")
	root.insertValue("E")
	root.insertValue("B")
	root.insertValue("K")
	root.insertValue("S")
	root.insertValue("Q")
	root.insertValue("M")

	// Add F.
	root.insertValue("F")

	fmt.Println("Sorted? ", root.isSorted())
	// Display the values in sorted order.
	fmt.Printf("Sorted values: %s\n", root.right.inorder())
	/*
		root.right.left = &Node{"X", nil, nil}
		fmt.Println("Sorted (+X)? ", root.isSorted())
		// Display the values in sorted order.
		fmt.Printf("Sorted values: %s\n", root.right.inorder())
	*/

	// Let the user search for values.
	for {
		// Get the target value.
		target := ""
		fmt.Printf("String: ")
		fmt.Scanln(&target)
		if len(target) == 0 {
			break
		}

		// Find the value's node.
		node := root.findValue(target)
		if node == nil {
			fmt.Printf("%s not found\n", target)
		} else {
			fmt.Printf("Found value %s\n", target)
		}
	}
}
