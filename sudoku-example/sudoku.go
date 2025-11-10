package sudokuexample

// SudokuState represents an immutable 9x9 Sudoku board
// 0 means empty cell
type SudokuState struct {
	Board [9][9]int
}

// PlaceNumber returns a new SudokuState with number placed at (row, col)
// Law I - Immutable operation
func (s SudokuState) PlaceNumber(row, col, num int) SudokuState {
	if row < 0 || row > 8 || col < 0 || col > 8 {
		return s // Invalid position, return unchanged
	}
	if num < 0 || num > 9 {
		return s // Invalid number, return unchanged
	}
	if s.Board[row][col] != 0 {
		return s // Cell already filled, return unchanged
	}

	// Create new board (immutable)
	newBoard := s.Board
	newBoard[row][col] = num

	return SudokuState{Board: newBoard}
}

// Merge combines two SudokuStates (associative operation for Law I)
// If both have same position filled with different numbers, keep the non-zero one
// If both have same number, keep it (idempotent)
// This allows distributed solving: two solvers work on different parts, then merge!
func (s SudokuState) Merge(other SudokuState) SudokuState {
	newBoard := s.Board

	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if s.Board[row][col] == 0 && other.Board[row][col] != 0 {
				// We have empty, other has value -> take other's value
				newBoard[row][col] = other.Board[row][col]
			} else if s.Board[row][col] != 0 && other.Board[row][col] == 0 {
				// We have value, other is empty -> keep ours
				newBoard[row][col] = s.Board[row][col]
			} else if s.Board[row][col] == other.Board[row][col] {
				// Both same (including both empty) -> keep it (idempotent)
				newBoard[row][col] = s.Board[row][col]
			} else {
				// Conflict: both have different non-zero values
				// Last-write-wins (or could flag conflict)
				// For true CRDT, this shouldn't happen if solvers coordinate
				newBoard[row][col] = other.Board[row][col]
			}
		}
	}

	return SudokuState{Board: newBoard}
}

// IsValid checks if current board state is valid (no conflicts)
func (s SudokuState) IsValid() bool {
	// Check rows
	for row := 0; row < 9; row++ {
		if !s.isValidSet(s.getRow(row)) {
			return false
		}
	}

	// Check columns
	for col := 0; col < 9; col++ {
		if !s.isValidSet(s.getCol(col)) {
			return false
		}
	}

	// Check 3x3 boxes
	for boxRow := 0; boxRow < 3; boxRow++ {
		for boxCol := 0; boxCol < 3; boxCol++ {
			if !s.isValidSet(s.getBox(boxRow, boxCol)) {
				return false
			}
		}
	}

	return true
}

// IsSolved checks if board is completely filled and valid
func (s SudokuState) IsSolved() bool {
	if !s.IsValid() {
		return false
	}

	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if s.Board[row][col] == 0 {
				return false
			}
		}
	}

	return true
}

// Helper methods

func (s SudokuState) getRow(row int) []int {
	return s.Board[row][:]
}

func (s SudokuState) getCol(col int) []int {
	result := make([]int, 9)
	for row := 0; row < 9; row++ {
		result[row] = s.Board[row][col]
	}
	return result
}

func (s SudokuState) getBox(boxRow, boxCol int) []int {
	result := make([]int, 9)
	idx := 0
	startRow := boxRow * 3
	startCol := boxCol * 3

	for row := startRow; row < startRow+3; row++ {
		for col := startCol; col < startCol+3; col++ {
			result[idx] = s.Board[row][col]
			idx++
		}
	}
	return result
}

func (s SudokuState) isValidSet(nums []int) bool {
	seen := make(map[int]bool)
	for _, num := range nums {
		if num == 0 {
			continue // Empty cells are OK
		}
		if seen[num] {
			return false // Duplicate!
		}
		seen[num] = true
	}
	return true
}

// CountFilled returns number of filled cells
func (s SudokuState) CountFilled() int {
	count := 0
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if s.Board[row][col] != 0 {
				count++
			}
		}
	}
	return count
}
