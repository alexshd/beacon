# gor-show

**"Show me how you think, and I will find my own way to do it"**

A pedagogical showcase demonstrating property-based testing with `lawtest` - a library that verifies mathematical properties of your code using group theory.

## What This Is

This is NOT a framework or production code. This is a **learning resource** that shows:

1. **The Problem** - Why traditional testing isn't enough
2. **The Solution** - How mathematical properties catch bugs tests miss
3. **The Practice** - Real examples you can run and modify

## Prerequisites

Install [Task](https://taskfile.dev) - a task runner for Go projects:

```bash
# macOS
brew install go-task

# Linux (snap)
sudo snap install task --classic

# Or download from https://taskfile.dev/installation
```

## Quick Start

```bash
# See all available commands
task --list

# Interactive demonstration
task show-me

# Run all tests in all examples
task test
```

## Project Structure

```
gor-show/
├── faulttest/              # Law I - Immutability demonstration
│   ├── state.go           # Immutable state implementation
│   ├── state_test.go      # Property tests with lawtest
│   └── README.md          # Detailed explanation
│
├── config-merge-example/   # lawtest usage example
│   ├── config.go          # Config merge with wrapper pattern
│   ├── config_test.go     # Traditional unit tests (pass ✓)
│   ├── config_fuzz_test.go # Fuzz tests
│   ├── config_law_test.go # Property tests (find bugs ✗)
│   ├── doc.go             # Package documentation
│   ├── LAWTEST_USAGE.md   # When to use lawtest
│   └── Taskfile.yml       # Test commands
│
├── lawtest-gen/           # Code analysis tool
│   └── main.go            # Generates lawtest skeletons
│
└── lawtest-check/         # Interactive guide
    └── main.go            # Helps decide if lawtest fits
```

## Examples

### 1. Fault Injection (faulttest/)

Demonstrates why immutability is a mathematical necessity, not a style choice.

```bash
# Run the demonstration
task show-me

# Run tests directly
task test-faulttest
```

**Key Insight**: Go's `defer` prevents deadlock but CANNOT prevent state corruption with mutable shared memory.

### 2. Config Merge (config-merge-example/)

Shows three testing approaches on the same code:

```bash
cd config-merge-example

# Traditional tests - PASS (but miss bugs)
task test-normal

# Property tests with lawtest - FAIL (catch bugs)
task test-law

# Fuzz tests
task test-fuzz

# Run all three
task
```

**Key Insight**: `DeepMerge` passes unit tests but fails associativity - a real bug caught by lawtest.

## Tools

### lawtest-gen

Analyzes Go code and generates lawtest skeleton:

```bash
cd lawtest-gen
go build
./lawtest-gen ../config-merge-example/config.go
```

Generates test file with:

- Function signature analysis
- Comparability checks
- Test skeletons with TODOs
- Wrapper pattern guidance

### lawtest-check

Interactive tool to determine if lawtest fits your use case:

```bash
cd lawtest-check
go build
./lawtest-check
```

Asks questions and provides guidance on whether lawtest is appropriate.

## Core Concepts

### What is lawtest?

A **specialized** testing library that verifies mathematical properties:

- **Associativity**: `(a ∘ b) ∘ c = a ∘ (b ∘ c)`
- **Immutability**: Operations don't mutate inputs
- **Parallel Safety**: Safe for concurrent execution

### When to Use lawtest

✅ **Good fit:**

- Binary operations: `func(T, T) T`
- Config merging, state aggregation
- Pure functions without side effects
- Concurrent-safe operations

❌ **Poor fit:**

- Operations with side effects (I/O, database)
- Context-dependent operations
- Business logic with complex rules
- Non-associative operations (that's OK!)

See `config-merge-example/LAWTEST_USAGE.md` for detailed guidelines.

## Learning Path

1. **Start here**: `task show-me` - Interactive demonstration
2. **Read**: `faulttest/README.md` - Why immutability matters
3. **Explore**: `config-merge-example/` - See lawtest in action
4. **Try**: Run `lawtest-gen` on your own code
5. **Decide**: Use `lawtest-check` for your use cases

## Development

```bash
# Install development tools
task install-tools

# Run all tests
task test

# Clean artifacts
task clean

# Tidy dependencies
task tidy
```

## Philosophy

This project embodies:

- **Show, don't tell** - Working examples over theoretical explanations
- **Guide, don't dictate** - Tools suggest, humans decide
- **Prove, don't assume** - Mathematical verification over intuition

## Resources

- [lawtest library](https://github.com/alexshd/lawtest)
- [Task documentation](https://taskfile.dev)
- Theory of Constraints applied to concurrent systems
- Group Theory in software verification

---

**Remember**: lawtest is a complement to traditional tests and fuzzing, not a replacement. Use it where mathematical properties matter.
