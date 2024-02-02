package loops_detection

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

// Get a last cell
func (list *LinkedList) getLastCell() *Cell {
	lastCell := list.sentinel
	for lastCell.next != nil {
		lastCell = lastCell.next
	}
	return lastCell
}

// Add a cell after me.
func (me *Cell) addAfter(after *Cell) {
	after.next = me.next
	me.next = after
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

func (list *LinkedList) toStringMax(separator string, max int) string {
	count := 0
	result := ""
	for cell := list.sentinel.next; cell != nil; cell = cell.next {
		result += cell.data
		count++
		if count >= max {
			break
		}
		if cell.next != nil {
			result += separator
		}
	}
	return result
}

func (list *LinkedList) isEmpty() bool {
	return list.sentinel.next == nil
}

func (list *LinkedList) hasLoop() bool {
	if list.isEmpty() || list.sentinel.next.next == nil {
		return false
	}

	tortoise := list.sentinel.next
	hare := list.sentinel.next.next

	for tortoise != hare && hare.next != nil && hare.next.next != nil && tortoise.next != nil {
		tortoise = tortoise.next
		hare = hare.next.next
	}
	return tortoise == hare
}

func LoopDetectionRun() {
	// Make a list from an array of values.
	values := []string{
		"0", "1", "2", "3", "4", "5",
	}
	list := makeLinkedList()
	list.addRange(values)

	fmt.Println(list.toString(" "))
	if list.hasLoop() {
		fmt.Println("Has loop")
	} else {
		fmt.Println("No loop")
	}
	fmt.Println()

	// Make cell 5 point to cell 2.
	list.sentinel.next.next.next.next.next.next = list.sentinel.next.next

	fmt.Println(list.toStringMax(" ", 10))
	if list.hasLoop() {
		fmt.Println("Has loop")
	} else {
		fmt.Println("No loop")
	}
	fmt.Println()

	// Make cell 4 point to cell 2.
	list.sentinel.next.next.next.next.next = list.sentinel.next.next

	fmt.Println(list.toStringMax(" ", 10))
	if list.hasLoop() {
		fmt.Println("Has loop")
	} else {
		fmt.Println("No loop")
	}
}
