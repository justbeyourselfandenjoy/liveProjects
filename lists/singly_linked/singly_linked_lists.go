package singlylinked

import "fmt"

type Cell struct {
	data string
	next *Cell
}

type LinkedList struct {
	sentinel *Cell
}

func makeLinkedList() LinkedList {
	list := LinkedList{}
	list.sentinel = &Cell{"SENTINEL", nil}
	return list
}

// Add a cell after me.
func (me *Cell) addAfter(after *Cell) {
	after.next = me.next
	me.next = after
}

// Get a last cell
func (list *LinkedList) getLastCell() *Cell {
	lastCell := list.sentinel
	for lastCell.next != nil {
		lastCell = lastCell.next
	}
	return lastCell
}

// Delete a cell after me.
func (me *Cell) deleteAfter() {
	if me.next == nil {
		panic("No a cell to delete")
	}
	me.next = me.next.next
}

func (list *LinkedList) addRange(values []string) {
	lastCell := list.getLastCell()
	if lastCell == nil {
		fmt.Println("addRange: lastCell is nil. Exiting")
		return
	}

	var cell *Cell
	for _, v := range values {
		cell = new(Cell)
		cell.data = v
		lastCell.addAfter(cell)
		lastCell = cell
	}
}

func (list *LinkedList) toString(separator string) string {
	result := ""
	for cell := list.sentinel.next; cell != nil; cell = cell.next {
		result += cell.data
		if cell.next != nil {
			result += separator
		}
	}
	return result
}

func (list *LinkedList) length() (len int) {
	len = 0
	for cell := list.sentinel.next; cell != nil; cell = cell.next {
		len++
	}
	return
}

func (list *LinkedList) isEmpty() bool {
	return list.sentinel.next == nil
}

func (list *LinkedList) push(value string) {
	lastCell := list.getLastCell()
	lastCell.addAfter(&Cell{data: value})
}

func (list *LinkedList) pop() (popped string) {
	popped = ""
	if list.sentinel.next == nil {
		return
	}
	cell := list.sentinel
	popped = cell.next.data
	for cell.next.next != nil {
		popped = cell.next.next.data
		cell = cell.next
	}
	cell.deleteAfter()
	return
}

func SinglyLinkedRun() {
	// Make a list from an array of values.
	greekLetters := []string{
		"α", "β", "γ", "δ", "ε",
	}
	list := makeLinkedList()
	list.addRange(greekLetters)
	fmt.Println(list.toString(" "))

	// Demonstrate a stack.
	stack := makeLinkedList()
	stack.push("Apple")
	stack.push("Banana")
	stack.push("Coconut")
	stack.push("Date")
	for !stack.isEmpty() {
		fmt.Printf("Popped: %-7s   Remaining %d: %s\n",
			stack.pop(),
			stack.length(),
			stack.toString(" "))
	}
}
