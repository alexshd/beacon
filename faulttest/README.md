# Fault Test - Fault Injection Testing

## What It Does

**Breaks things on purpose to prove your isolation works.**

This package:

1. **Injects failures** - Deliberately crashes code with panics
2. **Verifies isolation** - Proves one failure doesn't corrupt everything
3. **Tests immutability** - Shows data can't be accidentally changed
4. **Demonstrates the problem** - Proves why Go's normal approach fails

## The Problem

In the Theory of Constraints applied to concurrent systems, Joe Armstrong identified that shared-memory concurrency creates a coupling point where:

- **$1 < r < 3$**: Stable, predictable behavior
- **$r > 3$**: Chaotic, geometric failure propagation

The `CriticalState` structure represents this vulnerability:

```go
type CriticalState struct {
    Config map[string]string  // Shared mutable state
    Lock   sync.Mutex          // Ultimate single point of failure
}
```

When a goroutine panics while holding the lock or after partial writes:

1. **Without `defer`**: Deadlock (system halts)
2. **With `defer`**: Lock released, but **partial writes persist** (state corruption)

### The Test Results

Running `go test -v ./internal/pivt/` demonstrates:

#### ✅ **Immutability Laws Enforced**

- `ImmutableOp`: State.Merge does not mutate inputs
- `ParallelSafe`: No race conditions detected
- `NoMutation`: Original states remain unchanged

#### ❌ **Shared Memory Violation Detected**

```
Geometric Failure Detected: State was partially corrupted and persisted.
Found corrupted state: map[key1:value_corrupt_PARTIAL]
```

**This proves the central thesis**: Go's `defer` prevents deadlock but **cannot prevent state corruption** in shared-memory models.

### The Solution: Functional State Management

The `State` type enforces immutability:

```go
type State struct {
    data map[string]string  // Private, immutable
}

func (s *State) Set(key, value string) *State {
    // Returns NEW state, never mutates original
    newData := make(map[string]string, len(s.data)+1)
    for k, v := range s.data {
        newData[k] = v
    }
    newData[key] = value
    return &State{data: newData}
}
```

## Mathematical Verification

Using `github.com/alexshd/lawtest`, we verify:

### Group Theory Properties

- **Associativity**: `(a ∘ b) ∘ c = a ∘ (b ∘ c)` ✅
- **Identity**: `empty ∘ x = x` ✅
- **Immutability**: No input mutation ✅
- **Parallel Safety**: No race conditions ✅

### Containment Properties

- **Panic Recovery**: Supervisor catches failures ✅
- **Lock Release**: System doesn't deadlock ✅
- **State Isolation**: Corrupted writes don't persist (ONLY with functional approach) ✅

## Implications

This package provides **executable proof** that:

1. **Go's primitives alone are insufficient** for r < 3 stability
2. **Functional constraints must be encoded** in the type system
3. **Mathematical laws** (verified by lawtest) provide the missing safety guarantee
4. **Abstract Algebra** is not theoretical - it's the practical solution to concurrency

## Running the Tests

```bash
# Run all fault injection tests
go test -v ./faulttest/

# Run with race detector
go test -race ./faulttest/

# Run benchmarks to measure overhead
go test -bench=. ./faulttest/
```

## Next Steps

1. **Implement Actor Model**: Use `State` as message type
2. **Add Supervision Trees**: Process restart with clean state
3. **Enforce at Compile Time**: Make mutation impossible, not just tested
4. **Extend lawtest**: Add more group-theoretic properties

## References

- Armstrong, J. "Making reliable distributed systems in the presence of software errors"
- Theory of Constraints (TOC) applied to concurrent systems
- Abstract Algebra for software verification
- `github.com/alexshd/lawtest` - Property-based testing using group theory
