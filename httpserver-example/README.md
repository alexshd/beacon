# HTTP Server Example - Law I + Law II

Demonstrates **visible isolation** with concurrent HTTP server using:

- **Law I**: Immutable TodoState operations (lawtest verified)
- **Law II**: Supervisor pattern with worker restart on panic

## Quick Start

```bash
# Run tests (Law I property verification)
go test -v

# Build and run server
go build -o httpserver ./cmd/main.go
./httpserver
```

Server runs on `http://localhost:8080`

## Endpoints

### GET /

View current todos

```bash
curl http://localhost:8080/
```

### POST /add

Add todo (with optional chaos injection)

```bash
# Normal add
curl -X POST http://localhost:8080/add \
  -H "Content-Type: application/json" \
  -d '{"title": "Learn Law I"}'

# Add with 50% failure probability (chaos mode)
curl -X POST http://localhost:8080/add \
  -H "Content-Type: application/json" \
  -d '{"title": "Test isolation", "inject_fault": 50}'
```

### GET /metrics

System health dashboard

```bash
curl http://localhost:8080/metrics
```

Shows:

- `requests_processed`: Total requests handled
- `failures_handled`: Panics caught and contained
- `worker_restarts`: Worker restarts by supervisor
- `r_eff`: Effective coupling (1 < r < 3 = stable)
- `status`: Stable | Warning | Chaos

### GET /status

Worker status and restart log

```bash
curl http://localhost:8080/status
```

## Isolation Demo

**Prove isolation works under chaos:**

```bash
# Terminal 1: Run server
./httpserver

# Terminal 2: Send concurrent requests with failures
for i in {1..100}; do
  curl -X POST http://localhost:8080/add \
    -H "Content-Type: application/json" \
    -d "{\"title\": \"Todo $i\", \"inject_fault\": 30}" &
done
wait

# Check metrics - all panics contained
curl http://localhost:8080/metrics

# Check todos - no corruption despite 30 panics
curl http://localhost:8080/
```

**Expected results:**

- ✅ All failures reported in `failures_handled`
- ✅ Workers restarted by supervisor (Law II)
- ✅ Todo state remains consistent (Law I)
- ✅ `r_eff` stays in stable zone (< 3.0)
- ✅ No request hangs or corrupts others

## What You're Seeing

**Law I (Immutability)**

- Every `/add` creates NEW TodoState
- Original state never mutated
- Verified by lawtest: `ImmutableOp`, `Associative`, `ParallelSafe`

**Law II (Supervisor)**

- Worker panics caught by `recover()`
- Supervisor notified via channel
- Worker "restarted" (tracked in metrics)
- System continues processing

**Isolation Proof**

- Panic in worker A doesn't affect worker B
- State before panic is unchanged
- Concurrent requests process during failures
- `r_eff` metric shows system remains stable

## Architecture

```
HTTP Request → Server.ProcessRequest()
                 ↓
               Select Worker (random)
                 ↓
         [Panic Recovery Wrapper]  ← Law II
                 ↓
         Read Current State (immutable)
                 ↓
         Apply Operation (new state)  ← Law I
                 ↓
         Update State (atomic)
                 ↓
         Return or Notify Supervisor
```

## Tests

Run property-based tests:

```bash
go test -v
```

Tests verify:

- ✅ `TestMergeImmutability` - Operations don't mutate
- ✅ `TestMergeAssociativity` - (a∘b)∘c = a∘(b∘c)
- ✅ `TestMergeParallelSafe` - Safe under concurrency

All use `lawtest` with custom equality for non-comparable TodoState.
