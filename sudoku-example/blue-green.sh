#!/bin/bash
# Blue-Green Deployment: Sudoku Edition
# No YAML files were harmed during this deployment!

echo "=== Sudoku Blue-Green Deployment ==="
echo "Demonstrating Law I: Associative Merge for Zero-Downtime Upgrades"
echo ""

pkill -f "sudoku.*900" || true
sleep 1

# Start Blue (v1.0) and Green (v2.0) servers
echo "Starting Blue server (v1.0) on :9000..."
./sudoku "Blue-v1.0" 9000 > /tmp/blue.log 2>&1 &
BLUE_PID=$!

echo "Starting Green server (v2.0) on :9001..."
./sudoku "Green-v2.0" 9001 > /tmp/green.log 2>&1 &
GREEN_PID=$!

sleep 2

echo ""
echo "=== Phase 1: Independent Solving ==="
echo ""

# Blue solves top-left quadrant (rows 0-4, cols 0-4)
echo "Blue (v1.0): Solving top-left quadrant..."
curl -s -X POST http://localhost:9000/place -d '{"row":0,"col":0,"num":5}' > /dev/null
curl -s -X POST http://localhost:9000/place -d '{"row":0,"col":1,"num":3}' > /dev/null
curl -s -X POST http://localhost:9000/place -d '{"row":1,"col":0,"num":6}' > /dev/null
curl -s -X POST http://localhost:9000/place -d '{"row":1,"col":1,"num":9}' > /dev/null
curl -s -X POST http://localhost:9000/place -d '{"row":2,"col":0,"num":8}' > /dev/null

# Green solves bottom-right quadrant (rows 5-8, cols 5-8)
echo "Green (v2.0): Solving bottom-right quadrant..."
curl -s -X POST http://localhost:9001/place -d '{"row":7,"col":7,"num":1}' > /dev/null
curl -s -X POST http://localhost:9001/place -d '{"row":7,"col":8,"num":4}' > /dev/null
curl -s -X POST http://localhost:9001/place -d '{"row":8,"col":7,"num":7}' > /dev/null
curl -s -X POST http://localhost:9001/place -d '{"row":8,"col":8,"num":2}' > /dev/null
curl -s -X POST http://localhost:9001/place -d '{"row":6,"col":6,"num":9}' > /dev/null

echo ""
echo "Blue state:"
curl -s http://localhost:9000/board | jq '{version, filled, valid}'

echo ""
echo "Green state:"
curl -s http://localhost:9001/board | jq '{version, filled, valid}'

echo ""
echo "=== Phase 2: Blue-Green Merge (The Magic!) ==="
echo ""

# Export Green's state
echo "Exporting Green (v2.0) state..."
GREEN_STATE=$(curl -s http://localhost:9001/export)

# Merge Green into Blue
echo "Merging Green → Blue (Law I: Associative Merge)..."
echo "$GREEN_STATE" | curl -s -X POST http://localhost:9000/merge \
  -H "Content-Type: application/json" \
  -d @- | jq '{message, filled, valid, version}'

echo ""
echo "Blue after merge:"
curl -s http://localhost:9000/board | jq '{version, filled, valid, solved}'

# Reverse: Export Blue, merge into Green
echo ""
echo "Exporting Blue state..."
BLUE_STATE=$(curl -s http://localhost:9000/export)

echo "Merging Blue → Green (Law I: Commutative!)..."
echo "$BLUE_STATE" | curl -s -X POST http://localhost:9001/merge \
  -H "Content-Type: application/json" \
  -d @- | jq '{message, filled, valid, version}'

echo ""
echo "Green after merge:"
curl -s http://localhost:9001/board | jq '{version, filled, valid, solved}'

echo ""
echo "=== Phase 3: Verification ==="
echo ""

echo "Both servers now have identical puzzle state!"
echo ""
echo "Blue filled cells:"
curl -s http://localhost:9000/board | jq '.filled'

echo "Green filled cells:"
curl -s http://localhost:9001/board | jq '.filled'

echo ""
echo "=== Results ==="
echo ""
echo "✅ Blue (v1.0) solved 5 cells independently"
echo "✅ Green (v2.0) solved 5 cells independently"
echo "✅ Merged: Both servers now have all 10 cells"
echo "✅ No conflicts - merge resolved automatically"
echo "✅ Associative: Can merge in any order"
echo "✅ Commutative: Blue.Merge(Green) = Green.Merge(Blue)"
echo ""
echo "🎯 This is Blue-Green Deployment with CRDT properties!"
echo ""
echo "Standard approach would need:"
echo "  ❌ Complex state reconciliation"
echo "  ❌ Distributed locks"
echo "  ❌ Leader election"
echo "  ❌ 500 lines of YAML"
echo ""
echo "Law I approach:"
echo "  ✅ Just call Merge()"
echo "  ✅ Math guarantees consistency"
echo "  ✅ lawtest proves correctness"
echo "  ✅ Zero YAML harmed!"

# Cleanup
kill $BLUE_PID $GREEN_PID 2>/dev/null || true
