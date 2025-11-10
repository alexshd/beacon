#!/bin/bash
# Quick test to see actual IDs generated

pkill -f "httpserver.*808" || true
sleep 1

./httpserver 8080 >/dev/null 2>&1 &
A_PID=$!
./httpserver 8081 >/dev/null 2>&1 &
B_PID=$!
sleep 2

echo "Adding 5 todos to Server A (should be IDs 10-14):"
for i in {1..5}; do
	curl -s -X POST http://localhost:8080/add -H "Content-Type: application/json" -d "{\"title\":\"A-$i\"}" >/dev/null
done
curl -s http://localhost:8080/ | jq '.todos[] | {id, title}'

echo ""
echo "Adding 5 todos to Server B (should be IDs 20-24):"
for i in {1..5}; do
	curl -s -X POST http://localhost:8081/add -H "Content-Type: application/json" -d "{\"title\":\"B-$i\"}" >/dev/null
done
curl -s http://localhost:8081/ | jq '.todos[] | {id, title}'

echo ""
echo "Merging B into A..."
curl -s http://localhost:8081/export | curl -s -X POST http://localhost:8080/merge -H "Content-Type: application/json" -d @- >/dev/null

echo ""
echo "Server A after merge (should have 10 todos with IDs 10-14 and 20-24):"
curl -s http://localhost:8080/ | jq '{count, next_id, ids: [.todos[] | .id]}'

kill $A_PID $B_PID 2>/dev/null || true
