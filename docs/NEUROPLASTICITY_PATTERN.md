# Neuroplasticity Pattern for Distributed Systems

## The Core Insight

Systems should adapt like brains: when one pathway fails, **discover alternate pathways that work**, **preserve that knowledge**, and **reroute traffic** accordingly.

Just as the brain reroutes visual cortex to process hearing when sight is lost, distributed systems can reroute traffic to services that have adapted to degraded conditions.

## The Problem with Traditional Systems

**Traditional Health Checks:**

```
Service healthy? → Route traffic ✅
Service unhealthy? → Stop routing ❌
```

**Problem:** Binary thinking. Services are either "up" or "down."

**Reality:** Services often work in **degraded modes**:

- DB slow? Use cache
- Auth down? Serve public reads
- Queue full? Buffer locally
- Network flaky? Batch requests

**Traditional systems throw away these working degraded states.**

---

## The Neuroplasticity Pattern

### Inspired by Brain Plasticity

**In Neurology:**

```
Visual cortex loses sight input
→ Brain detects hearing still works
→ Brain reroutes visual cortex to process audio
→ Enhanced hearing ability (compensation)
```

**In Distributed Systems:**

```
Service loses DB connection
→ Service detects cache still works
→ Supervisor learns service CAN work with cache
→ Reroute read-heavy traffic here (compensation)
```

### The Architecture

```
┌──────────────────────────────────────────┐
│  Service (Agent)                         │
│  ──────────────────────────────────────  │
│  • Reports working state continuously    │
│  • Discovers own capabilities            │
│  • Adapts to degraded conditions         │
│  • Self-aware of what it CAN do          │
└────────────────┬─────────────────────────┘
                 │ Working state reports
                 │ (not just "healthy/unhealthy")
                 ↓
┌──────────────────────────────────────────┐
│  Supervisor (Memory + Learning)          │
│  ──────────────────────────────────────  │
│  • Preserves what works                  │
│  • Learns service capabilities           │
│  • Remembers adaptation patterns         │
│  • Makes rerouting decisions             │
└────────────────┬─────────────────────────┘
                 │ Rerouting commands
                 │ (based on learned working states)
                 ↓
┌──────────────────────────────────────────┐
│  Router / Load Balancer                  │
│  ──────────────────────────────────────  │
│  • Routes based on capabilities          │
│  • Adapts topology dynamically           │
│  • Leverages degraded states             │
└──────────────────────────────────────────┘
```

---

## How It Works

### Phase 1: Normal Operation

```go
// Service reports healthy state
agent.Report(WorkingState{
    ServiceID: "api-server-1",
    Status:    "healthy",
    Dependencies: map[string]DepStatus{
        "db":    {Status: "up", Latency: 5 * time.Millisecond},
        "auth":  {Status: "up", Latency: 10 * time.Millisecond},
        "cache": {Status: "up", Latency: 1 * time.Millisecond},
    },
    Capabilities: []Capability{
        {Type: "read", Mode: "normal"},
        {Type: "write", Mode: "normal"},
    },
    Adaptations: []string{},
})

// Supervisor learns: "Service works in normal mode"
supervisor.Learn("api-server-1", "normal mode")
```

**Routing:** All traffic → Service (normal routing)

---

### Phase 2: Degradation Detected - Adaptation

```go
// DB gets slow, but service discovers it can still work!
agent.Report(WorkingState{
    ServiceID: "api-server-1",
    Status:    "degraded-but-working", // Key: STILL WORKING
    Dependencies: map[string]DepStatus{
        "db":    {Status: "SLOW", Latency: 500 * time.Millisecond}, // ❌
        "auth":  {Status: "up", Latency: 10 * time.Millisecond},
        "cache": {Status: "up", Latency: 1 * time.Millisecond},     // ✅ Compensating!
    },
    Capabilities: []Capability{
        {Type: "read", Mode: "cached"},          // Adapted!
        {Type: "write", Mode: "write-through"},  // Adapted!
    },
    Adaptations: []string{
        "using-cache-for-reads",
        "buffering-writes-to-cache",
        "async-db-sync",
    },
})

// Supervisor: "Whoa! Service found a way to work!"
supervisor.Learn("api-server-1", "cache-compensated mode")

// Neuroplasticity: Reroute based on learned capability
supervisor.Reroute(RoutingRule{
    Pattern: "read-heavy traffic",
    Target:  "api-server-1", // Still good for reads!
})
supervisor.Reroute(RoutingRule{
    Pattern: "write-heavy traffic",
    Target:  "api-server-2", // Use service with fast DB
})
```

**Routing:** Read traffic → Service (adapted), Write traffic → Other service

---

### Phase 3: Critical Failure - Autonomous Mode

```go
// Auth service dies completely, but service STILL serves!
agent.Report(WorkingState{
    ServiceID: "api-server-1",
    Status:    "autonomous", // Service found way to work without critical dep!
    Dependencies: map[string]DepStatus{
        "db":    {Status: "slow", Latency: 500 * time.Millisecond},
        "auth":  {Status: "DOWN"},  // ❌ Critical dependency DOWN
        "cache": {Status: "up"},     // ✅ Compensating
    },
    Capabilities: []Capability{
        {Type: "read", Mode: "public-cached"}, // No auth needed!
    },
    Adaptations: []string{
        "serving-public-data-from-cache",
        "no-auth-required-for-reads",
        "writes-disabled",
    },
    Degradation: "Cannot serve authenticated requests",
})

// Supervisor: "Service discovered autonomous operation!"
supervisor.Learn("api-server-1", "autonomous-cached mode")

// Neuroplasticity: Leverage what works, route around what doesn't
supervisor.Reroute(RoutingRule{
    Pattern: "public read requests",
    Target:  "api-server-1", // Can handle without auth!
})
supervisor.Reroute(RoutingRule{
    Pattern: "authenticated requests",
    Target:  "api-server-3", // Route to service with auth
})
supervisor.Reroute(RoutingRule{
    Pattern: "writes",
    Target:  "queue", // Queue for later
})
```

**Routing:** Public reads → Service (autonomous), Auth required → Other service, Writes → Queue

---

## Code Structure

### Working State Report

```go
type WorkingState struct {
    ServiceID    string
    Status       ServiceStatus // healthy | degraded-but-working | autonomous | down
    Dependencies map[string]DepStatus
    Capabilities []Capability
    Adaptations  []Adaptation // What the service learned to do!
    Degradation  string       // What it can't do
}

type DepStatus struct {
    Name    string
    Status  string        // up | slow | degraded | down
    Latency time.Duration
}

type Capability struct {
    Type string // read | write | auth | etc.
    Mode string // normal | cached | degraded | public-only
}

type Adaptation struct {
    Name        string // "using-cache-for-reads"
    Trigger     string // "db-slow"
    Fallback    string // "cache"
    Performance string // "acceptable" | "degraded"
}
```

### Supervisor Memory

```go
type SupervisorMemory struct {
    // What each service can do in different conditions
    LearnedModes map[string][]WorkingMode

    // Current routing topology
    Routes map[Pattern][]ServiceID

    // History of adaptations (for learning)
    AdaptationHistory []AdaptationEvent
}

type WorkingMode struct {
    ModeName     string
    Capabilities []Capability
    Conditions   map[string]string // Dependencies required
    Performance  PerformanceMetrics
}
```

### Supervisor Learning & Rerouting

```go
func (s *Supervisor) ProcessWorkingState(state WorkingState) {
    // 1. Learn what the service can do
    if state.Status != "down" {
        mode := WorkingMode{
            ModeName:     s.inferModeName(state),
            Capabilities: state.Capabilities,
            Conditions:   s.extractConditions(state.Dependencies),
        }
        s.Memory.Learn(state.ServiceID, mode)
    }

    // 2. Reroute traffic based on learned capabilities
    s.reroute(state)
}

func (s *Supervisor) reroute(state WorkingState) {
    // Clear old routes for this service
    s.Memory.ClearRoutes(state.ServiceID)

    // Add new routes based on current capabilities
    for _, capability := range state.Capabilities {
        pattern := s.capabilityToPattern(capability)
        s.Memory.Routes[pattern] = append(
            s.Memory.Routes[pattern],
            state.ServiceID,
        )
    }

    // Notify router of topology change
    s.Router.UpdateTopology(s.Memory.Routes)
}
```

---

## Real-World Example: E-commerce API

### Scenario: Database Failure During Black Friday

**Without Neuroplasticity:**

```
1. DB becomes overloaded
2. Health checks fail
3. All API servers marked unhealthy
4. Traffic stops
5. ❌ Site down during peak sales
```

**With Neuroplasticity:**

```
1. DB becomes overloaded (200ms → 2s latency)

2. API server detects and adapts:
   report: {
     status: "degraded-but-working",
     adaptations: ["using-redis-cache", "stale-reads-ok"],
     capabilities: ["read-products", "read-cart"],
   }

3. Supervisor learns: "Server can serve reads from cache"

4. Supervisor reroutes:
   - Product browsing → Cache-adapted servers ✅
   - Cart viewing → Cache-adapted servers ✅
   - Checkout → Servers with DB access ✅
   - Writes → Queue for async processing ✅

5. ✅ Site stays up with degraded performance
   - Most features work (browsing, viewing cart)
   - Critical features protected (checkout)
   - Writes buffered for later
```

**Result:** Revenue saved, customers served, system adapted.

---

## Comparison to Traditional Patterns

### Health Check Pattern

**Traditional:**

```go
func HealthCheck() bool {
    return db.Ping() && auth.Ping() && cache.Ping()
}
// Returns: true or false
// Routing: all or nothing
```

**Problem:** Throws away partial functionality

---

### Circuit Breaker Pattern

**Traditional:**

```go
if circuitBreaker.IsOpen("db") {
    return ErrorServiceUnavailable
}
```

**Problem:** Stops trying, doesn't discover alternatives

---

---

## The Pattern in Code

**Report what you CAN do, not just what's broken:**

```go
func ReportWorkingState() WorkingState {
    state := WorkingState{Status: "healthy"}

    // Try DB
    if err := db.Ping(); err != nil {
        state.Status = "degraded"
        state.Dependencies["db"] = DepStatus{Status: "down"}

        // But can we still work?
        if cache.Available() {
            state.Adaptations = append(state.Adaptations,
                "using-cache-instead-of-db")
            state.Capabilities = []Capability{
                {Type: "read", Mode: "cached"},
            }
        }
    }

    return state // Report capabilities, not just health
}
```

Services discover what they can do under degraded conditions. Supervisors learn and preserve those capabilities. Routers adapt topology accordingly.

---

---

## Connection to Group Theory

### Composable Working States

Working states from multiple services can be **merged associatively**:

```go
// Service A can do: read-cached, write-queue
stateA := WorkingState{
    Capabilities: []Capability{
        {Type: "read", Mode: "cached"},
        {Type: "write", Mode: "queued"},
    },
}

// Service B can do: read-db, write-db
stateB := WorkingState{
    Capabilities: []Capability{
        {Type: "read", Mode: "direct"},
        {Type: "write", Mode: "direct"},
    },
}

// Combined system capability
merged := stateA.Merge(stateB)
// Result: {read-cached, read-direct, write-queued, write-direct}

// Prove merge is associative!
lawtest.Associative(t, MergeWorkingState, genState, stateEqual)
```

**Why this matters:** Supervisor can understand **total system capability** by composing individual service states.

### Idempotent Learning

```go
// Learning the same mode multiple times = same result
supervisor.Learn(serviceID, mode)
supervisor.Learn(serviceID, mode) // Idempotent
// Result: Same learned state

lawtest.Idempotent(t, LearnMode, genMode, modeEqual)
```

**Why this matters:** Safe to re-learn, handles duplicate reports.

---

## Benefits

### System-Level

✅ **Graceful Degradation** - Services work in reduced capacity instead of failing  
✅ **Self-Healing** - System adapts without manual intervention  
✅ **Resilience** - Partial failures don't cascade  
✅ **Dynamic Topology** - Routing adapts to current capabilities  
✅ **Cost Optimization** - Use degraded services instead of spinning up replicas

### Operational

✅ **Reduced Alert Fatigue** - Services self-adapt to common issues  
✅ **Better Observability** - Know what services CAN do, not just what's broken  
✅ **Faster Recovery** - System routes around problems automatically  
✅ **Predictable Behavior** - Adaptation patterns can be learned and replayed

### Business

✅ **Higher Availability** - Services stay up in degraded mode  
✅ **Revenue Protection** - Critical flows work even when deps fail  
✅ **Customer Experience** - Gradual degradation vs hard failure

---

## Implementation Guidelines

### For Services (Agents)

1. **Report working state continuously** (not just health checks)
2. **Discover what you CAN do** when dependencies fail
3. **Be self-aware** of capabilities and limitations
4. **Adapt gracefully** - find alternate pathways

```go
// Good: Report capabilities
report := WorkingState{
    Status: "degraded-but-working",
    Capabilities: []Capability{
        {Type: "read", Mode: "cached"},
    },
}

// Bad: Just fail
return errors.New("DB down")
```

### For Supervisors

1. **Learn from agents** - preserve working states
2. **Remember adaptations** - build knowledge base
3. **Reroute intelligently** - leverage what works
4. **Don't discard degraded states** - use them!

```go
// Good: Learn and reroute
supervisor.Learn(serviceID, state.Capabilities)
supervisor.Reroute(basedOn: capabilities)

// Bad: Binary routing
if healthy { route } else { reject }
```

### For Routers

1. **Route by capability** not just by "up/down"
2. **Adapt topology dynamically** as services adapt
3. **Prefer working degraded** over failing completely

---

## When to Use

### Perfect For

✅ **Critical systems** that must stay up  
✅ **Unpredictable failures** where services must adapt  
✅ **Cost-sensitive** systems that want to use degraded capacity  
✅ **Complex dependencies** where partial failures are common

### Not Suitable For

❌ **Simple systems** with few dependencies (overkill)  
❌ **Hard real-time** systems (adaptation has latency)  
❌ **Stateless services** (no adaptation needed)

---

## Trade-offs

### Costs

- **Complexity**: Services must report rich state
- **Overhead**: Continuous state reporting
- **Learning curve**: New mental model

### Benefits

- **Resilience**: System adapts automatically
- **Availability**: Partial functionality > no functionality
- **Self-healing**: Less manual intervention

---

## Summary

The Neuroplasticity Pattern brings **brain-inspired adaptation** to distributed systems:

1. **Services discover** what they can do under degraded conditions
2. **Supervisors learn** and preserve working states
3. **Routers adapt** topology based on capabilities
4. **System self-heals** by rerouting around failures

**Key Insight:** Don't throw away degraded states. **Learn from them, preserve them, leverage them.**

Just as the brain reroutes neural pathways when one sense fails, distributed systems can reroute traffic when one dependency fails.

**"If it still works (even degraded), use it!"**

---

## References

- Neuroplasticity: How brains adapt to injury
- Graceful degradation in distributed systems
- Self-healing systems
- Adaptive routing algorithms
- CRDT composition (for working state merging)

---

## Future Work

- **Machine learning** on adaptation patterns
- **Predictive rerouting** based on historical adaptations
- **Automated capability discovery** through exploration
- **Cross-service learning** - share adaptation strategies

---

_"The brain doesn't give up when one pathway fails. Neither should your distributed system."_
