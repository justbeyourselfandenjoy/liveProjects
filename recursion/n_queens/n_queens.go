package nqueens

import (
	"fmt"
	"time"
)

func makeBoard(size int) [][]string {
	board := make([][]string, size)
	for r := range board {
		board[r] = make([]string, size)
		for c := range board[r] {
			board[r][c] = "."
		}
	}
	return board
}

func dumpBoard(board [][]string) {
	for r := range board {
		for c := range board[r] {
			fmt.Printf("%s ", board[r][c])
		}
		fmt.Println()
	}
}

// Return true if this series of squares contains at most one queen.
func seriesIsLegal(board [][]string, r0, c0, dr, dc int) bool {
	numRows := len(board)
	numCols := numRows
	hasQueen := false

	for r, c := r0, c0; r < numRows && c < numCols && r >= 0 && c >= 0; r, c = r+dr, c+dc {
		if board[r][c] == "Q" {
			if hasQueen {
				return false
			}
			hasQueen = true
		}
	}
	return true
}

// Return true if the board is legal.
func boardIsLegal(board [][]string) bool {
	numRows := len(board)

	for r := 0; r < numRows; r++ {
		if !seriesIsLegal(board, r, 0, 0, 1) {
			return false
		}
		if !seriesIsLegal(board, r, 0, 1, 1) {
			return false
		}
		if !seriesIsLegal(board, r, numRows-1, 1, -1) {
			return false
		}
	}

	for c := 0; c < numRows; c++ {
		if !seriesIsLegal(board, 0, c, 1, 0) {
			return false
		}
		if !seriesIsLegal(board, 0, c, 1, 1) {
			return false
		}
		if !seriesIsLegal(board, 0, c, 1, -1) {
			return false
		}
	}
	return true
}

// Return true if the board is legal and a solution.
func boardIsASolution(board [][]string) bool {
	if boardIsLegal(board) {
		boardSize := len(board)
		numQueens := 0
		for r := range board {
			for c := range board[r] {
				if board[r][c] == "Q" {
					numQueens++
				}
			}
		}
		return boardSize == numQueens
	}
	return false
}

// Try placing a queen at position [r][c].
// Return true if we find a legal board.
func placeQueens1(board [][]string, numRows, r, c, numPlaced int) bool {
	if r >= numRows || c >= numRows || numPlaced == numRows {
		return boardIsASolution(board)
	}
	// Find the next square.
	nextR := r
	nextC := c + 1
	if nextC >= numRows {
		nextR += 1
		nextC = 0
	}
	if placeQueens1(board, numRows, nextR, nextC, numPlaced) {
		return true
	}
	board[r][c] = "Q"
	if placeQueens1(board, numRows, nextR, nextC, numPlaced+1) {
		return true
	}
	board[r][c] = "."
	return false
}

// Try to place a queen in this column.
// Return true if we find a legal board.
func placeQueens4(board [][]string, numRows, c int) bool {
	if c == numRows {
		if boardIsLegal(board) {
			return true
		}
		return false
	}
	if !boardIsLegal(board) {
		return false
	}
	for r := 0; r < numRows; r++ {
		board[r][c] = "Q"
		if placeQueens4(board, numRows, c+1) {
			return true
		}
		board[r][c] = "."
	}
	return false
}

func NQueensRun() {
	fmt.Println("Running NQueensRun()")

	const numRows = 8
	board := makeBoard(numRows)

	start := time.Now()
	//	success := placeQueens1(board, numRows, 0, 0, 0)
	//success := placeQueens2(board, numRows, 0, 0, 0)
	//success := placeQueens3(board, numRows, 0, 0, 0)
	success := placeQueens4(board, numRows, 0)

	elapsed := time.Since(start)
	if success {
		fmt.Println("Success!")
		dumpBoard(board)
	} else {
		fmt.Println("No solution")
	}
	fmt.Printf("Elapsed: %f seconds\n", elapsed.Seconds())
}
