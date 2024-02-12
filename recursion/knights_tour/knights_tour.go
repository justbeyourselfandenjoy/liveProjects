package knightstour

import (
	"fmt"
	"time"
)

// The board dimensions.
const numRows = 8
const numCols = numRows

// Whether we want an open or closed tour.
const requireClosedTour = false

// Value to represent a square that we have not visited.
const unvisited = -1

// Define offsets for the knight's movement.
type Offset struct {
	dr, dc int
}

var moveOffsets []Offset

var numCalls int64

func initializeOffsets() {
	moveOffsets = []Offset{
		{-2, +1},
		{-1, +2},
		{+1, +2},
		{+2, +1},
		{+2, -1},
		{+1, -2},
		{-1, -2},
		{-2, -1},
	}
}

func makeBoard(numRows, numCols int) [][]int {
	board := make([][]int, numRows)
	for r := range board {
		board[r] = make([]int, numCols)
		for c := range board[r] {
			board[r][c] = unvisited
		}
	}
	return board
}

func dumpBoard(board [][]int) {
	for r := range board {
		for c := range board[r] {
			fmt.Printf("%02d ", board[r][c])
		}
		fmt.Println()
	}
}

func moveIsOutOfTheBoard(offset Offset, numRows, numCols, curRow, curCol int) bool {
	if curRow+offset.dr >= numRows || curRow+offset.dr < 0 || curCol+offset.dc >= numCols || curCol+offset.dc < 0 {
		return true
	}
	return false
}

func squareIsVisited(visitsNum int) bool {
	return visitsNum != unvisited
}

// Try to extend a knight's tour starting at (curRow, curCol).
// Return true or false to indicate whether we have found a solution.
func findTour(board [][]int, numRows, numCols, curRow, curCol, numVisited int) bool {
	numCalls++
	var offset Offset
	if numVisited == numRows*numCols {
		if !requireClosedTour {
			return true
		} else {
			for _, offset = range moveOffsets {
				if moveIsOutOfTheBoard(offset, numRows, numCols, curRow, curCol) {
					continue
				}
				if board[curRow+offset.dr][curCol+offset.dc] == 0 { //?
					return true
				}

			}
			return false
		}
	}
	//	fmt.Println("numVisited !== numRows*numCols", numVisited, numRows, numCols)
	for _, offset = range moveOffsets {
		r := curRow + offset.dr
		c := curCol + offset.dc
		//		fmt.Printf("Move offset {%v,%v} to {%v,%v} \n", offset.dr, offset.dc, r, c)
		if moveIsOutOfTheBoard(offset, numRows, numCols, curRow, curCol) {
			//			fmt.Println("Illegal (out of the board)")
			continue
		}
		if squareIsVisited(board[r][c]) {
			//			fmt.Println("Illegal (visited)", board[r][c])
			continue
		}
		board[r][c] = numVisited
		if findTour(board, numRows, numCols, r, c, numVisited+1) {
			return true
		}
		board[r][c] = unvisited
	}
	return false
}

func KnightsTourRun() {
	numCalls = 0

	// Initialize the move offsets.
	initializeOffsets()

	// Create the blank board.
	board := makeBoard(numRows, numCols)

	// Try to find a tour.
	start := time.Now()
	board[0][0] = 0
	if findTour(board, numRows, numCols, 0, 0, 1) {
		fmt.Println("Success!")
	} else {
		fmt.Println("Could not find a tour.")
	}
	elapsed := time.Since(start)
	dumpBoard(board)
	fmt.Printf("%f seconds\n", elapsed.Seconds())
	fmt.Printf("%d calls\n", numCalls)
}
