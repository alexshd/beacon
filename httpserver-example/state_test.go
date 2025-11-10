package httpserver

import (
	"reflect"
	"testing"
	"time"

	"github.com/alexshd/lawtest"
)

// TodoStateWrapper wraps TodoState to make it comparable (pointer wrapper pattern)
type TodoStateWrapper struct {
	state *TodoState
}

func WrapMerge(a, b *TodoStateWrapper) *TodoStateWrapper {
	merged := a.state.Merge(*b.state)
	return &TodoStateWrapper{state: &merged}
}

// Custom equality for TodoStateWrapper
func todoStateEqual(a, b *TodoStateWrapper) bool {
	return reflect.DeepEqual(a.state, b.state)
}

// Test that Merge operation is immutable
func TestMergeImmutability(t *testing.T) {
	gen := func() *TodoStateWrapper {
		state := TodoState{
			Todos: []Todo{
				{ID: 1, Title: "Test", Completed: false, CreatedAt: time.Now()},
			},
			NextID: 2,
		}
		return &TodoStateWrapper{state: &state}
	}

	lawtest.ImmutableOpCustom(t, WrapMerge, gen, todoStateEqual)
}

// Test that Merge is associative
func TestMergeAssociativity(t *testing.T) {
	counter := 0
	gen := func() *TodoStateWrapper {
		counter++
		state := TodoState{
			Todos: []Todo{
				{ID: counter, Title: "Todo", Completed: false, CreatedAt: time.Now()},
			},
			NextID: counter + 1,
		}
		return &TodoStateWrapper{state: &state}
	}

	lawtest.AssociativeCustom(t, WrapMerge, gen, todoStateEqual)
}

// Test parallel safety
func TestMergeParallelSafe(t *testing.T) {
	counter := 0
	gen := func() *TodoStateWrapper {
		counter++
		state := TodoState{
			Todos: []Todo{
				{ID: counter, Title: "Test", Completed: false, CreatedAt: time.Now()},
			},
			NextID: counter + 1,
		}
		return &TodoStateWrapper{state: &state}
	}

	lawtest.ParallelSafeCustom(t, WrapMerge, gen, todoStateEqual, 100)
}
