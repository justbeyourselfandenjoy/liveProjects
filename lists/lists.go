package main

import (
	doublylinked "justbeyourselfandenjoy/lists/doubly_linked"
	loops_detection "justbeyourselfandenjoy/lists/loops"
	singlylinked "justbeyourselfandenjoy/lists/singly_linked"
)

func main() {
	singlylinked.SinglyLinkedRun()
	loops_detection.LoopDetectionRun()
	doublylinked.DoublyLinkedRun()
}
