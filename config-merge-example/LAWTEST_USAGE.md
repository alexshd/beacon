# lawtest Usage Definition

## What lawtest IS

A **property-based testing library** that verifies **mathematical properties** of operations using **group theory**.

It tests whether your code follows fundamental laws like:

- Associativity: `(a ∘ b) ∘ c = a ∘ (b ∘ c)`
- Immutability: Operations don't mutate inputs
- Parallel Safety: Safe for concurrent execution
- Identity: `empty ∘ x = x`

## What lawtest IS NOT

- NOT a replacement for unit tests
- NOT a replacement for fuzz testing
- NOT a general-purpose property testing framework
- NOT for testing business logic or complex workflows

## When lawtest WORKS (Use It)

### ✅ Binary Operations

Operations that combine two things of the same type:

```go
func Merge(a, b Config) Config
func Add(a, b Number) Number
func Combine(a, b State) State
```

### ✅ Data Structures with Operations

Types that have merge/combine/append operations:

- Configuration merging
- State aggregation
- Event reduction
- Collection operations (but see constraints below)

### ✅ Comparable Types

Types that Go can compare with `==`:

- Primitives: `int`, `string`, `bool`
- Structs of comparables
- **Pointers** (even if they point to non-comparables)

### ✅ Immutable Operations

Functions that should NOT mutate inputs:

- Pure functions
- Functional transformations
- Concurrent-safe operations

## When lawtest DOES NOT WORK (Don't Use It)

### ❌ Non-Comparable Types (Without Wrapper)

```go
type Config map[string]interface{}  // ❌ Maps not comparable
type Cache []Item                   // ❌ Slices not comparable
type Func func()                    // ❌ Functions not comparable
```

**Solution**: Wrap in a pointer-based struct (see ConfigWrapper example)

### ❌ Operations That Aren't Associative

Many real-world operations are NOT associative:

```go
// String concatenation with separators
Join("/", a, b)  // (a/b)/c ≠ a/(b/c)

// Subtraction
a - b - c ≠ a - (b - c)

// Division
a / b / c ≠ a / (b / c)
```

Don't force lawtest on these - they'll fail and that's OK.

### ❌ Operations with Side Effects

```go
func Save(a, b Data) Data {
    db.Write(a)  // ❌ Side effect
    db.Write(b)  // ❌ Side effect
    return merge(a, b)
}
```

lawtest assumes pure operations. Side effects break immutability testing.

### ❌ Context-Dependent Operations

```go
func Merge(a, b Config, env Environment) Config
```

lawtest works with binary operations `(a, b) -> c`, not operations that depend on external state.

### ❌ Operations on Different Types

```go
func Apply(config Config, override Override) Config
```

lawtest requires `(T, T) -> T` (same type in, same type out).

## Decision Tree

```
Is it a binary operation (a, b) -> c?
├─ NO → Don't use lawtest
└─ YES
   │
   Is the type comparable OR can you wrap it?
   ├─ NO → Don't use lawtest
   └─ YES
      │
      Should the operation be associative?
      ├─ NO → Don't use lawtest (or expect failures)
      └─ YES
         │
         Should the operation be immutable?
         ├─ NO → Don't use ImmutableOp test
         └─ YES → ✅ USE LAWTEST
```

## Checklist: Can I Use lawtest?

```
[ ] My operation has signature: (T, T) -> T
[ ] Type T is comparable OR I can wrap it with pointers
[ ] Operation should be associative (order doesn't matter)
[ ] Operation should be immutable (no mutation)
[ ] Operation is pure (no side effects)
[ ] I want to verify mathematical properties, not specific outputs
```

If ALL checkboxes are YES → Use lawtest

If ANY checkbox is NO → Consider traditional tests or fuzz tests instead

## Common Use Cases

### ✅ Good Fit

- Config merging (wrap map in pointer struct)
- Event sourcing (reduce events to state)
- CRDT operations (conflict-free replicated data)
- State machines (transition composition)
- Message aggregation
- Functional state updates

### ❌ Poor Fit

- HTTP handlers (side effects, context)
- Database operations (side effects, external state)
- Business logic (complex rules, not associative)
- String formatting (not associative)
- Mathematical operations with special cases (division, overflow)

## Example: Should I Use lawtest?

### Scenario 1: Shopping Cart

```go
func AddItem(cart Cart, item Item) Cart
```

- Binary? NO (takes different types)
- Comparable? Maybe
- **Decision: NO** - Use regular tests

### Scenario 2: Config Merge

```go
func Merge(a, b Config) Config
```

- Binary? YES ✓
- Comparable? NO, but can wrap ✓
- Associative? Should be ✓
- Immutable? Should be ✓
- **Decision: YES** - Perfect for lawtest

### Scenario 3: Balance Update

```go
func UpdateBalance(current, delta int) int
```

- Binary? YES ✓
- Comparable? YES ✓
- Associative? Addition is ✓
- Immutable? Ints are immutable ✓
- **Decision: YES** - Good for lawtest

### Scenario 4: File Path Join

```go
func JoinPath(a, b string) string
```

- Binary? YES ✓
- Comparable? YES ✓
- Associative? NO - path separators break it ✗
- **Decision: NO** - Will fail, not a bug

## Summary

**lawtest is specialized for:**

- Verifying algebraic properties
- Catching subtle concurrency bugs
- Proving operations are safe to parallelize

**lawtest is NOT for:**

- General correctness testing
- Business logic validation
- Integration testing
- Anything with side effects

Use it as a **complement** to traditional tests and fuzzing, not a replacement.
