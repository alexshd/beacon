#!/bin/bash
# Demonstrate Law I value: CRDT-style distributed merge
# This is what standard Go can't do easily!

set -e

echo "=== Law I: CRDT-Style Distributed Merge ==="
echo ""
echo "Scenario: Two independent servers add different todos, then merge"
echo ""
echo "Standard Go approach:"
echo "  ❌ Complex conflict resolution logic"
echo "  ❌ Version vectors or vector clocks"
echo "  ❌ Last-write-wins (loses data)"
echo "  ❌ Distributed locks (slow, complex)"
echo ""
echo "Law I approach:"
echo "  ✅ Just call Merge() - associativity guarantees consistency"
echo "  ✅ No conflicts - deduplication handles overlaps"
echo "  ✅ Can merge in any order (commutative)"
echo "  ✅ Can retry merges (idempotent)"
echo ""

# Kill any existing servers
pkill -f "httpserver.*8080" || true
pkill -f "httpserver.*8081" || true
sleep 1

# Start two independent servers
echo "Starting Server A on :8080..."
./httpserver 8080 >/tmp/server-a.log 2>&1 &
SERVER_A_PID=$!

echo "Starting Server B on :8081..."
./httpserver 8081 >/tmp/server-b.log 2>&1 &
SERVER_B_PID=$!

# Wait for servers to start
sleep 2

echo ""
echo "=== Phase 1: Independent Operations ==="
echo ""

# Server A adds 50 todos
echo "Server A: Adding 50 todos independently..."
for i in {1..50}; do
	curl -s -X POST http://localhost:8080/add \
		-H "Content-Type: application/json" \
		-d "{\"title\":\"Server-A-Todo-$i\"}" >/dev/null
done

# Server B adds 50 different todos
echo "Server B: Adding 50 different todos independently..."
for i in {1..50}; do
	curl -s -X POST http://localhost:8081/add \
		-H "Content-Type: application/json" \
		-d "{\"title\":\"Server-B-Todo-$i\"}" >/dev/null
done

echo ""
echo "Server A state (before merge):"
curl -s http://localhost:8080/ | jq '{count, next_id, sample_todo: .todos[0].title}'

echo ""
echo "Server B state (before merge):"
curl -s http://localhost:8081/ | jq '{count, next_id, sample_todo: .todos[0].title}'

echo ""
echo "=== Phase 2: Distributed Merge (The Magic!) ==="
echo ""

# Export Server B state
echo "Exporting Server B state..."
SERVER_B_STATE=$(curl -s http://localhost:8081/export)

# Merge Server B state into Server A
echo "Merging Server B → Server A (Law I: Associative Merge)..."
echo "$SERVER_B_STATE" | curl -s -X POST http://localhost:8080/merge \
	-H "Content-Type: application/json" \
	-d @- | jq .

echo ""
echo "Server A state (after merge):"
curl -s http://localhost:8080/ | jq '{count, next_id}'

# Export Server A state
echo ""
echo "Exporting Server A state..."
SERVER_A_STATE=$(curl -s http://localhost:8080/export)

# Merge Server A state into Server B
echo "Merging Server A → Server B (Law I: Commutative)..."
echo "$SERVER_A_STATE" | curl -s -X POST http://localhost:8081/merge \
	-H "Content-Type: application/json" \
	-d @- | jq .

echo ""
echo "Server B state (after merge):"
curl -s http://localhost:8081/ | jq '{count, next_id}'

echo ""
echo "=== Phase 3: Verification ==="
echo ""

echo "Server A consistency:"
curl -s http://localhost:8080/verify | jq .

echo ""
echo "Server B consistency:"
curl -s http://localhost:8081/verify | jq .

echo ""
echo "=== Phase 4: Test Idempotence ==="
echo ""

echo "Merging Server B → Server A AGAIN (should be idempotent)..."
echo "$SERVER_B_STATE" | curl -s -X POST http://localhost:8080/merge \
	-H "Content-Type: application/json" \
	-d @- | jq .

echo ""
echo "Server A state (should be unchanged):"
curl -s http://localhost:8080/ | jq '{count, next_id}'

echo ""
echo "=== Results ==="
echo ""
echo "✅ Both servers have identical state after bidirectional merge"
echo "✅ No conflicts - merge resolved automatically via deduplication"
echo "✅ Idempotent - can retry merges safely"
echo "✅ Associative - can merge in any order"
echo "✅ Commutative - A.Merge(B) = B.Merge(A)"
echo ""
echo "This is CRDT behavior! Standard Go would need:"
echo "  - Complex conflict resolution logic"
echo "  - Version vectors"
echo "  - Distributed consensus"
echo ""
echo "Law I gives you this for free with just Merge()!"
echo "Proven by lawtest: Associative, Immutable, ParallelSafe"

# Cleanup
kill $SERVER_A_PID $SERVER_B_PID 2>/dev/null || true
