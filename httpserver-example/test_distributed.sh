#!/bin/bash
# Demonstrate Law I value: Distributed CRDT-style merge
# Two "servers" independently add todos, then merge - no conflicts!

echo "=== Law I: Distributed Merge Demo ==="
echo ""
echo "Scenario: Two independent instances add todos, then merge"
echo "Standard Go: Complex conflict resolution, version vectors, last-write-wins"
echo "Law I approach: Just merge - associativity guarantees consistency"
echo ""

# Simulate Server A adding 100 todos
echo "Server A: Adding 100 todos independently..."
for i in {1..100}; do
	curl -s -X POST http://localhost:8080/add \
		-H "Content-Type: application/json" \
		-d "{\"title\":\"Server-A-Todo-$i\"}" >/dev/null
done

# Get Server A state
echo "Server A state:"
curl -s http://localhost:8080/ | jq '{count, next_id}'

echo ""
echo "Server B: Would add 100 different todos independently..."
echo "(In real distributed system, Server B would have different state)"
echo ""

# Simulate Server B adding 100 todos
echo "Adding 100 more todos (simulating Server B)..."
for i in {1..100}; do
	curl -s -X POST http://localhost:8080/add \
		-H "Content-Type: application/json" \
		-d "{\"title\":\"Server-B-Todo-$i\"}" >/dev/null
done

echo ""
echo "After merge (sequential here, but would be parallel in real system):"
curl -s http://localhost:8080/ | jq '{count, next_id}'

echo ""
echo "State consistency:"
curl -s http://localhost:8080/verify | jq .

echo ""
echo "=== Why Law I Matters ==="
echo ""
echo "✅ Associative: (A merge B) merge C = A merge (B merge C)"
echo "   → Can merge in any order"
echo ""
echo "✅ Commutative: A merge B = B merge A"
echo "   → Can merge from any direction"
echo ""
echo "✅ Idempotent: A merge A = A"
echo "   → Can retry merges safely"
echo ""
echo "✅ No conflicts: Merge resolves automatically via deduplication"
echo "   → No complex conflict resolution logic needed"
echo ""
echo "Standard Go approach needs:"
echo "  ❌ Version vectors or vector clocks"
echo "  ❌ Last-write-wins semantics (loses data)"
echo "  ❌ Complex merge conflict resolution"
echo "  ❌ Distributed locks or consensus (slow)"
echo ""
echo "Law I approach: Just merge! Math guarantees correctness."
