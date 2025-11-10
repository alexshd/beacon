package sudokuexample

import (
	"reflect"
	"testing"

	"github.com/alexshd/lawtest"
)

// Wrapper for lawtest (makes it comparable via pointer)
type SudokuStateWrapper struct {
	state *SudokuState
}

func WrapMerge(a, b *SudokuStateWrapper) *SudokuStateWrapper {
	merged := a.state.Merge(*b.state)
	return &SudokuStateWrapper{state: &merged}
}

func sudokuEqual(a, b *SudokuStateWrapper) bool {
	return reflect.DeepEqual(a.state, b.state)
}

// Test that Merge operation is immutable
func TestMergeImmutability(t *testing.T) {
	gen := func() *SudokuStateWrapper {
		state := SudokuState{}
		state = state.PlaceNumber(0, 0, 5)
		state = state.PlaceNumber(1, 1, 3)
		return &SudokuStateWrapper{state: &state}
	}

	lawtest.ImmutableOpCustom(t, WrapMerge, gen, sudokuEqual)
}

// Test that Merge is associative
// (A merge B) merge C = A merge (B merge C)
func TestMergeAssociativity(t *testing.T) {
	counter := 0
	gen := func() *SudokuStateWrapper {
		state := SudokuState{}
		// Each generator creates a board with a unique cell filled
		row := counter / 9
		col := counter % 9
		state = state.PlaceNumber(row, col, (counter%9)+1)
		counter++
		return &SudokuStateWrapper{state: &state}
	}

	lawtest.AssociativeCustom(t, WrapMerge, gen, sudokuEqual)
}

// Test parallel safety - multiple goroutines can merge simultaneously
func TestMergeParallelSafe(t *testing.T) {
	counter := 0
	gen := func() *SudokuStateWrapper {
		state := SudokuState{}
		row := counter / 9
		col := counter % 9
		state = state.PlaceNumber(row, col, (counter%9)+1)
		counter++
		return &SudokuStateWrapper{state: &state}
	}

	lawtest.ParallelSafeCustom(t, WrapMerge, gen, sudokuEqual, 100)
}

// Test PlaceNumber immutability
func TestPlaceNumberImmutability(t *testing.T) {
	original := SudokuState{}
	original = original.PlaceNumber(0, 0, 5)

	// Place another number
	modified := original.PlaceNumber(1, 1, 3)

	// Original should be unchanged
	if original.Board[1][1] != 0 {
		t.Errorf("PlaceNumber mutated original state! Expected 0, got %d", original.Board[1][1])
	}

	// Modified should have both numbers
	if modified.Board[0][0] != 5 || modified.Board[1][1] != 3 {
		t.Errorf("PlaceNumber didn't create correct new state")
	}
}

// Test merging partial solutions
func TestDistributedSolving(t *testing.T) {
	// Solver A works on top half
	solverA := SudokuState{}
	solverA = solverA.PlaceNumber(0, 0, 5)
	solverA = solverA.PlaceNumber(0, 1, 3)
	solverA = solverA.PlaceNumber(1, 0, 6)

	// Solver B works on bottom half
	solverB := SudokuState{}
	solverB = solverB.PlaceNumber(7, 7, 9)
	solverB = solverB.PlaceNumber(8, 8, 1)
	solverB = solverB.PlaceNumber(8, 7, 4)

	// Merge solutions (blue-green deployment!)
	merged := solverA.Merge(solverB)

	// Check all numbers are present
	if merged.Board[0][0] != 5 || merged.Board[0][1] != 3 {
		t.Errorf("Merge lost solver A's work")
	}
	if merged.Board[7][7] != 9 || merged.Board[8][8] != 1 {
		t.Errorf("Merge lost solver B's work")
	}

	// Should have 6 filled cells
	if merged.CountFilled() != 6 {
		t.Errorf("Expected 6 filled cells, got %d", merged.CountFilled())
	}

	// Test commutativity: A.Merge(B) = B.Merge(A)
	mergedReverse := solverB.Merge(solverA)
	if !reflect.DeepEqual(merged, mergedReverse) {
		t.Errorf("Merge is not commutative!")
	}
}

// Test idempotence: A.Merge(A) = A
func TestMergeIdempotence(t *testing.T) {
	state := SudokuState{}
	state = state.PlaceNumber(0, 0, 5)
	state = state.PlaceNumber(1, 1, 3)

	merged := state.Merge(state)

	if !reflect.DeepEqual(state, merged) {
		t.Errorf("Merge is not idempotent!")
	}
}
