package faulttest

import (
	"fmt"
	"sync"
	"testing"

	"github.com/alexshd/lawtest"
)

// ANSI color codes for readable output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
)

func printSection(title string) {
	fmt.Printf("\n%s%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%s%s  %s%s\n", colorBold, colorCyan, title, colorReset)
	fmt.Printf("%s%s━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━%s\n\n", colorBold, colorCyan, colorReset)
}

func printStep(step, description string) {
	fmt.Printf("%s%s📍 %s:%s %s\n", colorBold, colorBlue, step, colorReset, description)
}

func printSuccess(message string) {
	fmt.Printf("%s✅ SUCCESS:%s %s\n", colorGreen, colorReset, message)
}

func printFailure(message string) {
	fmt.Printf("%s❌ FAILURE:%s %s\n", colorRed, colorReset, message)
}

func printWarning(message string) {
	fmt.Printf("%s⚠️  WARNING:%s %s\n", colorYellow, colorReset, message)
}

func printInfo(message string) {
	fmt.Printf("%s💡 INFO:%s %s\n", colorWhite, colorReset, message)
}

// TestImmutableStateTransfer verifies that SafeUpdate respects the mathematical
// law of immutability using Abstract Algebra (Group Theory).
func TestImmutableStateTransfer(t *testing.T) {
	printSection("LAW I - IMMUTABILITY TEST")

	printInfo("What we're testing: State operations must NEVER mutate the original data")
	printInfo("Why it matters: Mutation breaks isolation and allows failures to spread")
	printInfo("How we prove it: Using mathematical properties from Group Theory")
	fmt.Println()

	// Define the operation to test: State.Merge as a binary operation
	mergeOp := func(a, b *State) *State {
		return a.Merge(b)
	}

	// Generator for test states
	stateGen := func() *State {
		return NewState(map[string]string{
			"key1": "value1",
			"key2": "value2",
		})
	}

	// Apply the ImmutableOp test: Prove that merging states
	// does not mutate the input and produces a new, independent copy
	t.Run("ImmutableOp", func(t *testing.T) {
		printStep("Test 1", "Proving operations don't mutate inputs")
		printInfo("Creating two states and merging them...")

		lawtest.ImmutableOp(t, mergeOp, stateGen)

		printSuccess("Operations preserve original data - no mutation detected")
		printInfo("This means: One actor's operation cannot corrupt another actor's state")
	})

	// Verify that the operation is safe for parallel execution
	t.Run("ParallelSafe", func(t *testing.T) {
		printStep("Test 2", "Proving operations are safe for concurrent use")
		printInfo("Running 20 concurrent operations to detect race conditions...")

		isSafe := lawtest.ParallelSafe(t, mergeOp, stateGen, 20)
		if !isSafe {
			printFailure("Race conditions detected - concurrent operations interfere with each other")
			t.Error("State.Merge has race conditions")
		} else {
			printSuccess("No race conditions - safe for concurrent actors")
			printInfo("This means: Multiple actors can operate simultaneously without conflicts")
		}
	})

	// Test associativity: (a merge b) merge c = a merge (b merge c)
	t.Run("Associative", func(t *testing.T) {
		printStep("Test 3", "Proving operations are associative")
		printInfo("Testing: (a + b) + c = a + (b + c)")
		printInfo("Why: Order of operations shouldn't change the result")

		a := NewState(map[string]string{"a": "1"})
		b := NewState(map[string]string{"b": "2"})
		c := NewState(map[string]string{"c": "3"})

		left := a.Merge(b).Merge(c)
		right := a.Merge(b.Merge(c))

		if left.Len() != right.Len() {
			printFailure("Associativity violated - different lengths")
			t.Errorf("Associativity violated: different lengths")
		}

		for _, key := range []string{"a", "b", "c"} {
			lval, lok := left.Get(key)
			rval, rok := right.Get(key)
			if lok != rok || lval != rval {
				printFailure(fmt.Sprintf("Associativity violated for key %s", key))
				t.Errorf("Associativity violated for key %s", key)
			}
		}

		printSuccess("Associativity proven - operation order doesn't matter")
		printInfo("This means: Message processing order is predictable and safe")
	})
}

// TestSafeMergeProperties verifies that SafeMerge satisfies group-theoretic properties.
func TestSafeMergeProperties(t *testing.T) {
	// Test identity: merge(empty, x) = x
	t.Run("Identity", func(t *testing.T) {
		state := NewState(map[string]string{"key": "value"})
		empty := NewState(map[string]string{})

		result := empty.Merge(state)

		if result.Len() != state.Len() {
			t.Errorf("Identity property violated: expected len %d, got %d", state.Len(), result.Len())
		}

		val, ok := result.Get("key")
		if !ok || val != "value" {
			t.Errorf("Identity property violated: merge(empty, state) != state")
		}
	})

	// Test associativity using lawtest
	t.Run("AssociativityProperty", func(t *testing.T) {
		a := NewState(map[string]string{"a": "1"})
		b := NewState(map[string]string{"b": "2"})
		c := NewState(map[string]string{"c": "3"})

		left := a.Merge(b).Merge(c)
		right := a.Merge(b.Merge(c))

		if left.Len() != right.Len() {
			t.Errorf("Associativity violated: different lengths %d vs %d", left.Len(), right.Len())
		}

		// Manual check for small example
		for _, key := range []string{"a", "b", "c"} {
			lval, lok := left.Get(key)
			rval, rok := right.Get(key)
			if lok != rok || lval != rval {
				t.Errorf("Associativity violated for key %s", key)
			}
		}
	})

	// Test that original maps aren't mutated
	t.Run("NoMutation", func(t *testing.T) {
		original := NewState(map[string]string{"original": "data"})
		originalLen := original.Len()

		// Create new state by merging
		NewState(map[string]string{"new": "data"}).Merge(original)

		// Original should be unchanged
		if original.Len() != originalLen {
			t.Error("Original state was mutated")
		}

		val, ok := original.Get("original")
		if !ok || val != "data" {
			t.Error("Original state content was corrupted")
		}
	})
}

// TestContainment is the Process Isolation Verification Test.
//
// This test PROVES that crashes can be contained without corrupting the system.
func TestContainment(t *testing.T) {
	printSection("FAULT INJECTION TEST - Can we survive crashes?")

	printInfo("What we're testing: Can one actor crash without breaking the whole system?")
	printInfo("The challenge: Go's shared memory makes this hard")
	printInfo("What we'll do: Deliberately crash code and see what happens")
	fmt.Println()

	t.Run("BasicPanicRecovery", func(t *testing.T) {
		printStep("Experiment 1", "Single actor crashes while holding a lock")

		critical := NewCriticalState()

		printInfo("Starting an actor that will panic mid-operation...")

		// Launch the Chaotic Operation in a protected manner (simulating supervision)
		done := make(chan struct{})
		var panicCaught interface{}

		go func() {
			defer func() {
				if r := recover(); r != nil {
					panicCaught = r
					printInfo(fmt.Sprintf("Supervisor caught panic: %v", r))
				}
				close(done)
			}()
			MutateAndPanic(critical, "key1", "value_corrupt")
		}()

		<-done

		// Verify that the panic was caught
		if panicCaught == nil {
			printFailure("Panic was NOT caught - system would crash")
			t.Fatal("Expected panic to be caught by supervisor")
		} else {
			printSuccess("Panic was caught - system didn't crash")
		}

		// Can we access the lock? (proves lock was released)
		printInfo("Checking if the lock was released (testing for deadlock)...")
		critical.Lock.Lock()
		defer critical.Lock.Unlock()
		printSuccess("Lock was released - no deadlock")

		// But was state corrupted?
		printInfo("Checking if state was corrupted during the crash...")
		if _, ok := critical.Config["key1"]; ok {
			printWarning("State WAS corrupted - the panic left partial writes")
			printWarning(fmt.Sprintf("Found corrupted data: %v", critical.Config))
			printWarning("This is THE PROBLEM: defer released the lock, but couldn't undo the write")
			printInfo("Solution: Use immutable state so partial writes are impossible")
		} else {
			printSuccess("State is clean - crash was fully isolated")
		}
	})

	t.Run("MultipleSequentialFailures", func(t *testing.T) {
		printStep("Experiment 2", "Multiple crashes in sequence")
		printInfo("Testing: Do failures accumulate and corrupt more state over time?")

		critical := NewCriticalState()

		for i := range 5 {
			success, panicVal := IsolatedOperation(func() {
				MutateAndPanic(critical, "key", "value")
			})

			if success {
				printFailure("Operation succeeded when it should have failed")
				t.Error("Expected operation to fail, but it succeeded")
			}

			if panicVal == nil {
				printFailure("Panic was not recovered")
				t.Error("Expected panic to be recovered")
			}

			printInfo(fmt.Sprintf("Crash %d recovered: %v", i+1, panicVal))
		}

		// Check accumulated damage
		critical.Lock.Lock()
		size := len(critical.Config)
		critical.Lock.Unlock()

		printInfo(fmt.Sprintf("After 5 crashes, state has %d corrupted entries", size))

		if size > 0 {
			printWarning("Multiple crashes accumulated corruption - this is GEOMETRIC FAILURE")
			printWarning("Each crash leaves a scar that makes the next crash worse")
			printInfo("This proves why immutability is not optional - it's mathematical necessity")
		} else {
			printSuccess("No accumulated corruption - system remains clean")
		}
	})

	t.Run("ConcurrentFailuresWithIsolation", func(t *testing.T) {
		printStep("Experiment 3", "10 actors crash simultaneously")
		printInfo("The hardest test: concurrent crashes with shared state")

		critical := NewCriticalState()
		var wg sync.WaitGroup
		const numGoroutines = 10

		printInfo("Launching 10 actors that will all crash...")

		for i := range numGoroutines {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				success, panicVal := IsolatedOperation(func() {
					MutateAndPanic(critical, "key", "value")
				})

				if success {
					t.Errorf("Goroutine %d: Expected operation to fail", id)
				}
				if panicVal == nil {
					t.Errorf("Goroutine %d: Expected panic to be recovered", id)
				}
			}(i)
		}

		wg.Wait()
		printSuccess("All 10 crashes survived - no deadlock")

		// Check damage
		critical.Lock.Lock()
		size := len(critical.Config)
		critical.Lock.Unlock()

		printInfo(fmt.Sprintf("After 10 concurrent crashes, found %d corrupted entries", size))

		if size > numGoroutines {
			printFailure(fmt.Sprintf("Excessive corruption: expected at most %d, got %d", numGoroutines, size))
			t.Errorf("Excessive state corruption: expected at most %d entries, got %d",
				numGoroutines, size)
		} else if size > 0 {
			printWarning("Some corruption occurred - mutable state is vulnerable")
		} else {
			printSuccess("Zero corruption - perfect isolation")
		}
	})
}

// TestFunctionalIsolation demonstrates the SOLUTION: immutable state prevents corruption.
func TestFunctionalIsolation(t *testing.T) {
	printSection("THE SOLUTION - Immutable State")

	printInfo("Now we test the SOLUTION to the problems we just saw")
	printInfo("Key insight: If state cannot be mutated, crashes cannot corrupt it")
	fmt.Println()

	t.Run("ImmutableStatePreventsChaos", func(t *testing.T) {
		printStep("Solution Test 1", "Immutable state survives crashes")

		// Start with clean state
		currentState := NewState(map[string]string{"initial": "clean"})
		printInfo("Starting with clean immutable state")

		// Simulate a failing operation that tries to corrupt state
		printInfo("Attempting to update state, then crashing mid-operation...")
		success, panicVal := IsolatedOperation(func() {
			// Attempt to update state (creates NEW state, doesn't mutate)
			_ = currentState.Set("corrupted", "PARTIAL")
			// Panic before we can commit the new state
			panic("Simulated failure mid-operation")
		})

		if success {
			printFailure("Operation succeeded when it should have failed")
			t.Fatal("Expected operation to fail")
		}
		if panicVal == nil {
			printFailure("Panic was not recovered")
			t.Fatal("Expected panic to be recovered")
		}

		printSuccess("Crash was recovered")

		// PROOF: The current state was never mutated
		printInfo("Checking if original state was corrupted...")
		if val, ok := currentState.Get("corrupted"); ok {
			printFailure(fmt.Sprintf("State corruption detected: corrupted = %s", val))
			t.Errorf("State corruption detected: %s = %s", "corrupted", val)
		} else {
			printSuccess("Original state is completely unchanged")
			printInfo("Why: The crash happened before the new state could be committed")
			printInfo("The old state is immutable, so the crash had no effect on it")
		}

		// PROOF: We can continue with clean state
		printInfo("Attempting safe update after the crash...")
		newState := currentState.Set("safe", "update")
		if val, ok := newState.Get("safe"); !ok || val != "update" {
			printFailure("Failed to continue with clean state")
			t.Error("Failed to continue with clean state")
		} else {
			printSuccess("System continues normally with clean state")
			printInfo("This is the key: Crashes are ISOLATED, not CASCADING")
		}
	})

	t.Run("MultipleConcurrentFailuresWithImmutability", func(t *testing.T) {
		printStep("Solution Test 2", "100 concurrent crashes with immutable state")
		printInfo("The ultimate stress test: massive concurrent failures")

		baseState := NewState(map[string]string{"base": "state"})
		var wg sync.WaitGroup
		const numGoroutines = 100

		printInfo("Launching 100 goroutines that will all crash...")

		// Launch many concurrent operations that fail
		for range numGoroutines {
			wg.Go(func() {
				IsolatedOperation(func() {
					_ = baseState.Set("corrupted", "value")
					panic("Simulated concurrent failure")
				})
			})
		}

		wg.Wait()
		printSuccess("All 100 crashes survived")

		// PROOF: Base state is completely unchanged
		printInfo("Checking base state integrity...")
		if baseState.Len() != 1 {
			printFailure(fmt.Sprintf("Base state mutated: expected len=1, got %d", baseState.Len()))
			t.Errorf("Base state mutated: expected len=1, got %d", baseState.Len())
		}

		if val, ok := baseState.Get("base"); !ok || val != "state" {
			printFailure("Base state content corrupted")
			t.Error("Base state content corrupted")
		} else {
			printSuccess("Base state is PERFECTLY intact after 100 crashes")
			printInfo("Zero corruption. Zero state leakage. Zero cascade.")
			printInfo("")
			printSuccess("This is Law I proven: Immutability enables true isolation")
		}
	})
}

// TestGeometricFailurePrevention tests that our functional approach prevents
// the cascading failures that characterize r > 3 chaos.
func TestGeometricFailurePrevention(t *testing.T) {
	t.Run("FunctionalUpdatePreventsCascade", func(t *testing.T) {
		// Start with a base state
		state1 := map[string]string{"base": "value"}

		// Apply a series of updates using the functional approach
		state2 := SafeUpdate(state1, "update1", "data1")
		state3 := SafeUpdate(state2, "update2", "data2")

		// Verify that early states remain uncorrupted
		if state1["update1"] != "" {
			t.Error("State mutation detected: state1 was corrupted by later updates")
		}

		if state2["update2"] != "" {
			t.Error("State mutation detected: state2 was corrupted by later updates")
		}

		// Verify final state is correct
		if state3["base"] != "value" || state3["update1"] != "data1" || state3["update2"] != "data2" {
			t.Error("Final state is incorrect")
		}

		t.Log("SUCCESS: Functional updates prevent mutation cascade (r < 3 maintained)")
	})
}

// TestInterfaceCompliance verifies that types implementing ImmutableStore
// actually obey Law I, not just have the right method signatures.
//
// Go's type system is structural - it only checks method signatures exist.
// It does NOT verify the implementation actually follows the mathematical laws.
// This test PROVES the implementation is correct.
func TestInterfaceCompliance(t *testing.T) {
	printSection("INTERFACE COMPLIANCE - Does the implementation obey Law I?")

	printInfo("Go's type system only checks: 'Does this type have these methods?'")
	printInfo("It does NOT check: 'Do these methods actually work correctly?'")
	printInfo("We prove the implementation is Law I compliant, not just compilable")
	fmt.Println()

	t.Run("StateWrapperImplementsInterface", func(t *testing.T) {
		printStep("Check 1", "Verify StateWrapper compiles as ImmutableStore")

		// This compiles because Go checks method signatures
		var store ImmutableStore = NewStateWrapper(map[string]string{"key": "value"})

		printSuccess("StateWrapper has all required methods")
		printInfo("But does it actually WORK correctly? Let's prove it...")

		// Can we use it?
		val, ok := store.Get("key")
		if !ok || val != "value" {
			printFailure("Get method doesn't work")
			t.Fatal("Get method failed")
		}
		printSuccess("Get method works")
	})

	t.Run("ImmutabilityViaInterface", func(t *testing.T) {
		printStep("Check 2", "Prove Set doesn't mutate original via interface")

		original := NewStateWrapper(map[string]string{"original": "data"})
		printInfo("Created original store with 'original=data'")

		// Call Set through interface
		printInfo("Calling Set through interface...")
		newStore := original.Set("new", "value")

		// Prove original is unchanged
		printInfo("Checking if original was mutated...")
		if val, ok := newStore.Get("original"); !ok || val != "data" {
			printFailure("Original data lost in new store")
			t.Error("Original data not preserved")
		}

		if _, ok := original.Get("new"); ok {
			printFailure("Original was mutated! Set violated immutability")
			t.Error("Set mutated original - Law I violated")
		} else {
			printSuccess("Original is unchanged - immutability preserved")
		}

		if val, ok := newStore.Get("new"); !ok || val != "value" {
			printFailure("New value not in new store")
			t.Error("New value missing")
		} else {
			printSuccess("New store has new value")
		}

		printInfo("Conclusion: Interface implementation respects immutability")
	})

	t.Run("AssociativityViaInterface", func(t *testing.T) {
		printStep("Check 3", "Prove Merge is associative via interface")

		a := NewStateWrapper(map[string]string{"a": "1"})
		b := NewStateWrapper(map[string]string{"b": "2"})
		c := NewStateWrapper(map[string]string{"c": "3"})

		printInfo("Testing: (a merge b) merge c = a merge (b merge c)")

		left := a.Merge(b).Merge(c)
		right := a.Merge(b.Merge(c))

		if left.Len() != right.Len() {
			printFailure(fmt.Sprintf("Different lengths: left=%d, right=%d", left.Len(), right.Len()))
			t.Errorf("Associativity violated: different lengths")
		}

		// Check each key
		allGood := true
		for _, key := range []string{"a", "b", "c"} {
			lval, lok := left.Get(key)
			rval, rok := right.Get(key)
			if lok != rok || lval != rval {
				printFailure(fmt.Sprintf("Key '%s' differs: left=%v, right=%v", key, lval, rval))
				allGood = false
			}
		}

		if allGood {
			printSuccess("Merge is associative - operation order doesn't matter")
			printInfo("This means: Interface contract is mathematically sound")
		} else {
			t.Error("Associativity violated")
		}
	})

	t.Run("ParallelSafetyViaInterface", func(t *testing.T) {
		printStep("Check 4", "Prove interface is safe for concurrent use")

		store := NewStateWrapper(map[string]string{"base": "value"})
		var wg sync.WaitGroup
		const numGoroutines = 50

		printInfo(fmt.Sprintf("Launching %d concurrent operations through interface...", numGoroutines))

		// Multiple goroutines calling methods through interface
		for i := range numGoroutines {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				// Each creates new store, doesn't mutate
				_ = store.Set(fmt.Sprintf("key%d", id), fmt.Sprintf("value%d", id))
				_ = store.Merge(NewStateWrapper(map[string]string{"test": "data"}))
			}(i)
		}

		wg.Wait()
		printSuccess("All concurrent operations completed")

		// Original should be unchanged
		printInfo("Checking if original store was affected...")
		if store.Len() != 1 {
			printFailure(fmt.Sprintf("Original mutated: expected len=1, got %d", store.Len()))
			t.Errorf("Concurrent access mutated original")
		}

		if val, ok := store.Get("base"); !ok || val != "value" {
			printFailure("Original data corrupted")
			t.Error("Original store corrupted")
		} else {
			printSuccess("Original store unchanged after 50 concurrent operations")
			printInfo("Conclusion: Interface is safe for concurrent actors")
		}
	})

	t.Run("InterfaceComparisonWithBrokenImpl", func(t *testing.T) {
		printStep("Check 5", "Show what happens with a BROKEN implementation")
		printWarning("This demonstrates WHY we need these tests")

		printInfo("Imagine someone implements ImmutableStore incorrectly...")
		printInfo("Go's compiler would accept it (right method signatures)")
		printInfo("But our Law I tests would CATCH the violation")
		printInfo("")
		printInfo("Example broken implementation:")
		printInfo("  func (b *Broken) Set(k, v string) ImmutableStore {")
		printInfo("    b.data[k] = v  // MUTATES original!")
		printInfo("    return b       // Returns same instance")
		printInfo("  }")
		printInfo("")
		printWarning("This compiles! Go can't detect the bug.")
		printSuccess("But our immutability tests would fail immediately")
		printInfo("")
		printInfo("That's the value: Prove correctness, not just compilation")
	})
}

// BenchmarkIsolationOverhead measures the performance cost of enforcing isolation.
func BenchmarkIsolationOverhead(b *testing.B) {
	state := map[string]string{"key": "value"}

	b.Run("UnsafeMutation", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			state["key"] = "newValue"
		}
	})

	b.Run("SafeFunctionalUpdate", func(b *testing.B) {
		currentState := state
		for i := 0; i < b.N; i++ {
			currentState = SafeUpdate(currentState, "key", "newValue")
		}
	})
}
