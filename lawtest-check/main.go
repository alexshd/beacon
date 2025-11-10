package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// lawtest-check - Interactive tool to determine if lawtest is appropriate for your use case

func main() {
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  lawtest Applicability Checker")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("This tool helps you decide if lawtest is appropriate for")
	fmt.Println("testing your operation.")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	score := 0
	total := 0

	questions := []struct {
		question string
		reason   string
		weight   int
	}{
		{
			"Does your operation have signature (T, T) -> T (same type in and out)?",
			"lawtest works with binary operations on a single type",
			10,
		},
		{
			"Is the type comparable (can use == in Go) OR can you wrap it with pointers?",
			"lawtest needs to compare values for equality checks",
			10,
		},
		{
			"Should the operation be associative? (a op b) op c = a op (b op c)",
			"Most lawtest value comes from verifying associativity",
			8,
		},
		{
			"Should the operation be immutable (not mutate inputs)?",
			"ImmutableOp test requires operations don't mutate",
			8,
		},
		{
			"Is the operation pure (no side effects like I/O, database, etc)?",
			"lawtest assumes pure operations for property testing",
			9,
		},
		{
			"Does operation order matter for correctness?",
			"If order matters, operation likely isn't associative",
			5,
		},
		{
			"Is this for concurrent/parallel code?",
			"lawtest excels at proving concurrent safety",
			6,
		},
	}

	for i, q := range questions {
		total += q.weight
		fmt.Printf("%d. %s\n", i+1, q.question)
		fmt.Printf("   Why: %s\n", q.reason)
		fmt.Print("   Answer (y/n): ")

		scanner.Scan()
		answer := strings.ToLower(strings.TrimSpace(scanner.Text()))

		if i == 5 { // "order matters" question - inverted logic
			if answer == "n" || answer == "no" {
				score += q.weight
			}
		} else {
			if answer == "y" || answer == "yes" {
				score += q.weight
			}
		}
		fmt.Println()
	}

	percentage := (float64(score) / float64(total)) * 100

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("  RESULT")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("\nScore: %d/%d (%.0f%%)\n\n", score, total, percentage)

	if percentage >= 80 {
		fmt.Println("✅ EXCELLENT FIT for lawtest")
		fmt.Println()
		fmt.Println("Your operation is a perfect candidate for property-based")
		fmt.Println("testing with lawtest. You should use:")
		fmt.Println("  • lawtest.ImmutableOp() - verify no mutation")
		fmt.Println("  • lawtest.Associative() - verify order independence")
		fmt.Println("  • lawtest.ParallelSafe() - verify concurrent safety")
		fmt.Println()
		fmt.Println("See config-merge-example for implementation patterns.")
	} else if percentage >= 60 {
		fmt.Println("⚠️  PARTIAL FIT for lawtest")
		fmt.Println()
		fmt.Println("lawtest can help, but with limitations:")
		fmt.Println("  • Some tests may fail (that's OK if property doesn't apply)")
		fmt.Println("  • You may need wrapper types for non-comparable types")
		fmt.Println("  • Consider using alongside traditional tests")
		fmt.Println()
		fmt.Println("Review LAWTEST_USAGE.md for decision guidance.")
	} else {
		fmt.Println("❌ POOR FIT for lawtest")
		fmt.Println()
		fmt.Println("lawtest is NOT recommended for this use case.")
		fmt.Println()
		fmt.Println("Better alternatives:")
		fmt.Println("  • Traditional unit tests - for specific examples")
		fmt.Println("  • Fuzz testing - for finding edge cases")
		fmt.Println("  • Integration tests - for side effects")
		fmt.Println()
		fmt.Println("lawtest works best with pure, associative, binary operations.")
	}

	fmt.Println()
	fmt.Println("See LAWTEST_USAGE.md for detailed guidelines.")
	fmt.Println()
}
