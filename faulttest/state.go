package faulttest

import (
	"fmt"
	"maps"
	"sync"
)

// State represents an immutable configuration state.
// We use a comparable wrapper to enable property-based testing with lawtest.
type State struct {
	data map[string]string
}

// NewState creates a new State with the given data.
func NewState(data map[string]string) *State {
	copied := make(map[string]string, len(data))
	for k, v := range data {
		copied[k] = v
	}
	return &State{data: copied}
}

// Get retrieves a value from the state.
func (s *State) Get(key string) (string, bool) {
	if s == nil || s.data == nil {
		return "", false
	}
	val, ok := s.data[key]
	return val, ok
}

// Set returns a new State with the key-value pair added.
// This enforces immutability - the original State is unchanged.
func (s *State) Set(key, value string) *State {
	newData := make(map[string]string, len(s.data)+1)
	for k, v := range s.data {
		newData[k] = v
	}
	newData[key] = value
	return &State{data: newData}
}

// Merge combines two states, with the other state's values taking precedence.
func (s *State) Merge(other *State) *State {
	if other == nil {
		return s
	}
	newData := make(map[string]string, len(s.data)+len(other.data))
	for k, v := range s.data {
		newData[k] = v
	}
	for k, v := range other.data {
		newData[k] = v
	}
	return &State{data: newData}
}

// Len returns the number of entries in the state.
func (s *State) Len() int {
	if s == nil || s.data == nil {
		return 0
	}
	return len(s.data)
}

// String implements fmt.Stringer for debugging.
func (s *State) String() string {
	if s == nil || s.data == nil {
		return "{}"
	}
	return fmt.Sprintf("%v", s.data)
}

// CriticalState represents a shared resource susceptible to Geometric System Failure.
// In conventional Go code, this structure embodies the coupling point (r) where
// concurrent access to shared memory can lead to catastrophic propagation of failures.
//
// This is the $r > 3$ vulnerability: a panic while holding the lock causes deadlock,
// and partial writes corrupt the system state.
type CriticalState struct {
	// Config holds the system's critical configuration.
	// Shared state susceptible to race conditions and partial-write corruption.
	Config map[string]string

	// Lock is the standard Go primitive that represents the 'ultimate single point of failure'.
	// If a goroutine panics while holding this lock, the entire system deadlocks.
	Lock sync.Mutex
}

// NewCriticalState creates a new vulnerable shared state.
func NewCriticalState() *CriticalState {
	return &CriticalState{
		Config: make(map[string]string),
	}
}

// MutateAndPanic simulates a chaotic operation that attempts to update shared state
// but panics while holding the lock or after partially writing.
//
// This function deliberately violates safety constraints to test whether our
// supervision and isolation mechanisms can contain the failure and prevent
// geometric propagation
//
//	(r to infty)
//
// The failure modes tested:
// 1. Lock acquisition followed by panic → deadlock (if not handled)
// 2. Partial write followed by panic → state corruption (if not isolated)
func MutateAndPanic(state *CriticalState, key, value string) {
	state.Lock.Lock() // Acquire lock (the 'ultimate single point of failure')
	defer state.Lock.Unlock()

	// Simulate complex, interruptible operation
	state.Config[key] = value + "_PARTIAL"

	// DELIBERATE PANIC: Simulating a logic error (division by zero, nil dereference)
	panic("Simulated unexpected failure: Logic error in critical section")

	// If the panic happens before Unlock, the system deadlocks (Geometric Failure).
	// If the panic happens after partial write, the state is corrupt.
	// The defer statement should save us, but the partial write still occurred.
}

// SafeUpdate is a functional approach to state updates that respects immutability.
// This function creates a new map rather than mutating the existing one,
// enforcing the mathematical law of Mandatory Isolation (Law I).
//
// By operating on immutable data structures, we suppress the coupling parameter
//
//	r to the stable zone (1 < r < 3).
func SafeUpdate(oldState map[string]string, key, value string) map[string]string {
	// Create a new map (immutability enforced)
	newState := make(map[string]string, len(oldState)+1)

	// Copy all existing entries
	maps.Copy(newState, oldState)

	// Apply the update to the new map
	newState[key] = value

	return newState
}

// SafeMerge combines two state maps functionally without mutation.
// This operation must satisfy the properties of a mathematical group:
// - Identity: merge(empty, x) = x
// - Associativity: merge(merge(a, b), c) = merge(a, merge(b, c))
func SafeMerge(state1, state2 map[string]string) map[string]string {
	result := make(map[string]string, len(state1)+len(state2))

	for k, v := range state1 {
		result[k] = v
	}

	for k, v := range state2 {
		result[k] = v // state2 overwrites state1 in case of key conflicts
	}

	return result
}

// IsolatedOperation wraps a potentially dangerous operation in a supervised context.
// This simulates Law II (Preemptive Supervision) by catching panics and preventing
// geometric failure propagation.
//
// Returns:
// - success: true if operation completed without panic
// - panicValue: the recovered panic value if one occurred
func IsolatedOperation(operation func()) (success bool, panicValue interface{}) {
	defer func() {
		if r := recover(); r != nil {
			success = false
			panicValue = r
		}
	}()

	operation()
	success = true
	return
}

// ImmutableStore defines the interface for Law I compliant storage.
// Any type implementing this interface MUST obey:
// 1. Immutability - operations never mutate the receiver
// 2. Associativity - operation order doesn't matter
// 3. Parallel Safety - safe for concurrent access
//
// Go's type system only checks method signatures exist.
// It does NOT verify the methods actually obey the mathematical laws.
// That's what our tests prove.
type ImmutableStore interface {
	// Get retrieves a value by key
	Get(key string) (string, bool)

	// Set returns a NEW store with the key-value pair added
	// The original store MUST remain unchanged (Law I)
	Set(key, value string) ImmutableStore

	// Merge combines two stores, returning a NEW store
	// MUST be associative: (a+b)+c = a+(b+c)
	Merge(other ImmutableStore) ImmutableStore

	// Len returns the number of entries
	Len() int
}

// StateWrapper wraps State to implement ImmutableStore interface
type StateWrapper struct {
	state *State
}

// NewStateWrapper creates a new wrapped state
func NewStateWrapper(data map[string]string) *StateWrapper {
	return &StateWrapper{state: NewState(data)}
}

func (sw *StateWrapper) Get(key string) (string, bool) {
	return sw.state.Get(key)
}

func (sw *StateWrapper) Set(key, value string) ImmutableStore {
	newState := sw.state.Set(key, value)
	return &StateWrapper{state: newState}
}

func (sw *StateWrapper) Merge(other ImmutableStore) ImmutableStore {
	otherWrapper, ok := other.(*StateWrapper)
	if !ok {
		return sw
	}
	merged := sw.state.Merge(otherWrapper.state)
	return &StateWrapper{state: merged}
}

func (sw *StateWrapper) Len() int {
	return sw.state.Len()
}

// Ensure StateWrapper implements ImmutableStore at compile time
var _ ImmutableStore = (*StateWrapper)(nil)
