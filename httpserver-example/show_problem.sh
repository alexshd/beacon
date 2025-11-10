#!/bin/bash
# Show the ID collision problem in naive distributed systems

echo "=== Problem: Naive ID generation causes collisions ==="
echo ""

# Kill any existing servers
pkill -f "httpserver.*8080" || true
pkill -f "httpserver.*8081" || true
sleep 1

# Start two servers
./httpserver 8080 >/dev/null 2>&1 &
SERVER_A_PID=$!
./httpserver 8081 >/dev/null 2>&1 &
SERVER_B_PID=$!
sleep 2

# Both add 5 todos
for i in {1..5}; do
	curl -s -X POST http://localhost:8080/add -H "Content-Type: application/json" -d "{\"title\":\"A-$i\"}" >/dev/null
	curl -s -X POST http://localhost:8081/add -H "Content-Type: application/json" -d "{\"title\":\"B-$i\"}" >/dev/null
done

echo "Server A todos (IDs will be 1,2,3,4,5):"
curl -s http://localhost:8080/ | jq '.todos[] | {id, title}'

echo ""
echo "Server B todos (IDs will ALSO be 1,2,3,4,5!):"
curl -s http://localhost:8081/ | jq '.todos[] | {id, title}'

echo ""
echo "Merging Server B into Server A..."
SERVER_B_STATE=$(curl -s http://localhost:8081/export)
echo "$SERVER_B_STATE" | curl -s -X POST http://localhost:8080/merge -H "Content-Type: application/json" -d @- >/dev/null

echo ""
echo "Server A after merge (PROBLEM: only 5 todos, not 10!):"
curl -s http://localhost:8080/ | jq '{count, todos: [.todos[] | {id, title}]}'

echo ""
echo "❌ ID collision: Both servers generated IDs 1-5"
echo "❌ Merge deduplication removed 'duplicate' IDs"
echo "❌ Lost 5 todos!"
echo ""
echo "This is WHY distributed systems need:"
echo "  - UUID/ULID instead of sequential IDs"
echo "  - Server-specific ID prefixes"
echo "  - Vector clocks"
echo "  - Or CRDT-specific ID schemes"

kill $SERVER_A_PID $SERVER_B_PID 2>/dev/null || true
