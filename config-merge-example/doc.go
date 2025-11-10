// Package configmerge demonstrates property-based testing with lawtest.
//
// This package shows how lawtest complements traditional testing and fuzzing
// by verifying mathematical properties that code should satisfy.
//
// # The Problem
//
// Configuration merging is common in applications. Developers write tests
// that check specific examples work correctly. These tests pass. Code ships.
// But subtle bugs exist that only appear in specific combinations.
//
// # Three Testing Approaches
//
// Traditional Tests (config_test.go):
//   - Test specific examples
//   - Pass for happy paths
//   - Miss edge cases and property violations
//
// Fuzz Tests (config_fuzz_test.go):
//   - Generate random inputs
//   - Find crashes and panics
//   - Slow to find logical bugs
//
// Property Tests (config_law_test.go):
//   - Test mathematical properties
//   - Fast discovery of violations
//   - Prove correctness, not just examples
//
// # What We Test
//
// Immutability: Operations should not mutate inputs
//   - Critical for concurrent systems
//   - Prevents action-at-a-distance bugs
//
// Associativity: (a + b) + c should equal a + (b + c)
//   - Required for parallel processing
//   - Order of operations shouldn't matter
//
// # The Discovery
//
// Normal tests pass. DeepMerge appears to work.
// But lawtest quickly discovers: DeepMerge is NOT associative.
//
// This is a real bug. It means:
//   - Results depend on merge order
//   - Cannot safely parallelize
//   - Cannot rely on the operation in distributed systems
//
// # Why This Matters
//
// Many "working" codebases have similar issues. Traditional tests give
// false confidence. Fuzzing might eventually find problems, but slowly.
//
// lawtest finds these bugs immediately by testing the properties
// the code must satisfy, not just whether it produces expected output
// for hand-picked inputs.
//
// # Running The Tests
//
//	go test -v                     # Run all tests
//	go test -v -run Normal         # Traditional tests only
//	go test -v -run Law            # Property tests only
//	go test -fuzz=FuzzMerge        # Fuzz testing
//
// # Watch Mode
//
//	task dev-configmerge           # Auto-run tests on file changes
package configmerge
