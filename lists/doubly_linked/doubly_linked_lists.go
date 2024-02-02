package doublylinked

import "fmt"

type Cell struct {
	data string
	next *Cell
	prev *Cell
}

type DoublyLinkedList struct {
	topSentinel    *Cell
	bottomSentinel *Cell
}

func makeDoublyLinkedList() DoublyLinkedList {
	list := DoublyLinkedList{}
	list.topSentinel = &Cell{"TOP_SENTINEL", nil, nil}
	list.bottomSentinel = &Cell{"BOTTOM_SENTINEL", nil, nil}
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

// Add a cell before me.
func (me *Cell) addBefore(before *Cell) {
	before.prev = me.prev
	before.next = me
	if me.prev != nil {
		me.prev.next = before
	}
	me.prev = before
}

// Delete a cell after me.
func (me *Cell) delete() {
	if me.data == "TOP_SENTINEL" || me.data == "BOTTOM_SENTINEL" {
		panic("Can't delete a sentinel")
	}
	me.prev.next = me.next
	me.next.prev = me.prev
}

func (me *Cell) deleteBefore() {
	if me.prev != nil {
		me.prev.delete()
	}
}

func (me *Cell) deleteAfter() {
	if me.next != nil {
		me.next.delete()
	}

}

func (list *DoublyLinkedList) addRange(values []string) {
	for _, v := range values {
		list.bottomSentinel.addBefore(&Cell{data: v, next: nil, prev: nil})
		//list.topSentinel.addAfter(&Cell{data: v, next: nil, prev: nil})
	}
}

func (list *DoublyLinkedList) toString(separator string) string {
	result := ""
	for cell := list.topSentinel.next; cell != list.bottomSentinel && cell != nil; cell = cell.next {
		result += cell.data
		if cell.next != list.bottomSentinel {
			result += separator
		}
	}
	return result
}

// Add an item to the top of the queue.
func (queue *DoublyLinkedList) enqueue(value string) {
	queue.topSentinel.addAfter(&Cell{data: value, next: nil, prev: nil})
}

// Remove an item from the bottom of the queue.
func (queue *DoublyLinkedList) dequeue() string {
	r := ""
	if queue.bottomSentinel.prev != nil {
		r = queue.bottomSentinel.prev.data
	}
	queue.bottomSentinel.deleteBefore()
	return r
}

// Add an item at the bottom of the deque.
func (deque *DoublyLinkedList) pushBottom(value string) {
	deque.bottomSentinel.addBefore(&Cell{data: value, prev: nil, next: nil})
}

// Add an item at the top of the deque.
func (deque *DoublyLinkedList) popBottom() string {
	return deque.dequeue()
}

// Add an item at the top of the deque.
func (deque *DoublyLinkedList) pushTop(value string) {
	deque.topSentinel.addAfter(&Cell{data: value, prev: nil, next: nil})
}

// Remove an item from the top of the deque.
func (deque *DoublyLinkedList) popTop() string {
	r := ""
	if deque.topSentinel.next != nil {
		r = deque.topSentinel.next.data
	}
	deque.topSentinel.deleteAfter()
	return r
}

func DoublyLinkedRun() {
	// Make a list from a slice of values.
	/*
		list := makeDoublyLinkedList()
		animals := []string{
			"Ant",
			"Bat",
			"Cat",
			"Dog",
			"Elk",
			"Fox",
		}
		list.addRange(animals)
		fmt.Println(list.toString(" "))
	*/

	// Test queue functions.
	fmt.Printf("*** Queue Functions ***\n")
	queue := makeDoublyLinkedList()
	queue.enqueue("Agate")
	queue.enqueue("Beryl")
	fmt.Printf("%s ", queue.dequeue())
	queue.enqueue("Citrine")
	fmt.Printf("%s ", queue.dequeue())
	fmt.Printf("%s ", queue.dequeue())
	queue.enqueue("Diamond")
	queue.enqueue("Emerald")
	for !queue.isEmpty() {
		fmt.Printf("%s ", queue.dequeue())
	}
	fmt.Printf("\n\n")

	// Test deque functions. Names starting
	// with F have a fast pass.
	fmt.Printf("*** Deque Functions ***\n")
	deque := makeDoublyLinkedList()
	deque.pushTop("Ann")
	deque.pushTop("Ben")
	fmt.Printf("%s ", deque.popBottom())
	deque.pushBottom("F-Cat")
	fmt.Printf("%s ", deque.popBottom())
	fmt.Printf("%s ", deque.popBottom())
	deque.pushBottom("F-Dan")
	deque.pushTop("Eva")
	for !deque.isEmpty() {
		fmt.Printf("%s ", deque.popBottom())
	}
	fmt.Printf("\n")
}
