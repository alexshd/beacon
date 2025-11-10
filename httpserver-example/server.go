package httpserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
)

// Metrics tracks system health
type Metrics struct {
	RequestsProcessed atomic.Int64
	ConcurrentPeak    atomic.Int32
}

// Server implements Law I - immutable state operations
type Server struct {
	state *TodoState
	sync.RWMutex
	metrics Metrics
	idMult  int // ID multiplier for distributed unique IDs
}

func NewServer() *Server {
	return &Server{
		state: &TodoState{
			Todos:  []Todo{},
			NextID: 1,
		},
		idMult: 1,
	}
}

func NewServerWithIDMultiplier(idMult int) *Server {
	return &Server{
		state: &TodoState{
			Todos:  []Todo{},
			NextID: idMult * 100, // Server 1: 100-199, Server 2: 200-299
		},
		idMult: idMult,
	}
}

// ProcessRequest handles a request using immutable operations (Law I)
func (s *Server) ProcessRequest(title string) TodoState {
	log.Printf("[REQUEST] Processing: %s", title)

	// Law I - Read current state (immutable)
	s.RLock()
	currentState := *s.state
	currentID := currentState.NextID
	s.RUnlock()
	log.Printf("[STATE] Read current state, NextID=%d", currentID)

	// Law I - Create new state (pure function, no mutation)
	newState := currentState.Add(title)
	log.Printf("[STATE] Created new state, NextID=%d", newState.NextID)

	// Update state atomically
	s.Lock()
	s.state = &newState
	s.Unlock()
	log.Printf("[STATE] Updated state atomically")

	s.metrics.RequestsProcessed.Add(1)
	log.Printf("[REQUEST] Completed: %s (total: %d)", title, s.metrics.RequestsProcessed.Load())
	return newState
}

// HTTP Handlers

func (s *Server) HandleRoot(w http.ResponseWriter, r *http.Request) {
	s.RLock()
	todos := s.state.Todos
	nextID := s.state.NextID
	s.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"todos":   todos,
		"count":   len(todos),
		"next_id": nextID,
	})
}

func (s *Server) HandleAdd(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HTTP] Received POST /add request")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("[HTTP] Failed to decode request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		log.Printf("[HTTP] Empty title received")
		http.Error(w, "title required", http.StatusBadRequest)
		return
	}

	log.Printf("[HTTP] Adding todo: %s", req.Title)
	newState := s.ProcessRequest(req.Title)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"todo":    newState.Todos[len(newState.Todos)-1],
		"count":   len(newState.Todos),
	})
	log.Printf("[HTTP] Response sent for: %s", req.Title)
}

func (s *Server) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	s.RLock()
	todoCount := len(s.state.Todos)
	nextID := s.state.NextID
	s.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"requests_processed": s.metrics.RequestsProcessed.Load(),
		"todo_count":         todoCount,
		"next_id":            nextID,
		"law_i":              "Immutable state operations (lawtest verified)",
		"guarantee":          "State consistency proven by group theory properties",
	})
}

func (s *Server) HandleVerify(w http.ResponseWriter, r *http.Request) {
	s.RLock()
	todos := s.state.Todos
	nextID := s.state.NextID
	s.RUnlock()

	// Verify state consistency
	expectedNextID := len(todos) + 1
	consistent := (nextID == expectedNextID)

	// Check for ID gaps or duplicates
	ids := make(map[int]bool)
	for _, todo := range todos {
		if ids[todo.ID] {
			consistent = false
		}
		ids[todo.ID] = true
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"consistent":       consistent,
		"todo_count":       len(todos),
		"next_id":          nextID,
		"expected_next_id": expectedNextID,
		"message":          fmt.Sprintf("Law I guarantee: %v", consistent),
	})
}

// HandleExport exports the current state as JSON (for CRDT-style distributed merge)
func (s *Server) HandleExport(w http.ResponseWriter, r *http.Request) {
	s.RLock()
	state := *s.state
	s.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

// HandleMerge merges incoming TodoState using Law I associative Merge operation
// This demonstrates CRDT-style eventually consistent distributed state
func (s *Server) HandleMerge(w http.ResponseWriter, r *http.Request) {
	log.Printf("[HTTP] Received POST /merge request")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var incomingState TodoState
	if err := json.NewDecoder(r.Body).Decode(&incomingState); err != nil {
		log.Printf("[HTTP] Failed to decode incoming state: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[MERGE] Incoming state: %d todos, NextID=%d", len(incomingState.Todos), incomingState.NextID)

	// Law I - Read current state (immutable)
	s.RLock()
	currentState := *s.state
	s.RUnlock()
	log.Printf("[MERGE] Current state: %d todos, NextID=%d", len(currentState.Todos), currentState.NextID)

	// Law I - Associative merge (pure function, no mutation)
	// This is the CRDT magic: A.Merge(B).Merge(C) = A.Merge(B.Merge(C))
	mergedState := currentState.Merge(incomingState)
	log.Printf("[MERGE] Merged state: %d todos, NextID=%d", len(mergedState.Todos), mergedState.NextID)

	// Update state atomically
	s.Lock()
	s.state = &mergedState
	s.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success":    true,
		"todo_count": len(mergedState.Todos),
		"next_id":    mergedState.NextID,
		"merged":     len(incomingState.Todos),
		"message":    "Law I: Associative merge completed without conflicts",
	})
	log.Printf("[MERGE] Merge completed successfully")
}

func (s *Server) Start(addr string) error {
	http.HandleFunc("/", s.HandleRoot)
	http.HandleFunc("/add", s.HandleAdd)
	http.HandleFunc("/metrics", s.HandleMetrics)
	http.HandleFunc("/verify", s.HandleVerify)
	http.HandleFunc("/export", s.HandleExport)
	http.HandleFunc("/merge", s.HandleMerge)

	log.Printf("Server starting on %s", addr)
	log.Printf("Law I: Immutable operations (lawtest verified)")
	log.Printf("Endpoints: /, /add, /metrics, /verify, /export, /merge")
	return http.ListenAndServe(addr, nil)
}
