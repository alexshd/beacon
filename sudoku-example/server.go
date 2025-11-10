package sudokuexample

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/alexshd/beacon/sudoku-example/web"
)

// Server implements Sudoku solver with Law I operations
type Server struct {
	state   *SudokuState
	stateMu sync.RWMutex
	version string
}

func NewServer(version string) *Server {
	return &Server{
		state:   &SudokuState{},
		version: version,
	}
}

// HandlePlace places a number on the board (immutable operation)
func (s *Server) HandlePlace(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Row int `json:"row"`
		Col int `json:"col"`
		Num int `json:"num"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[%s] Placing %d at (%d,%d)", s.version, req.Num, req.Row, req.Col)

	// Law I - Read current state (immutable)
	s.stateMu.RLock()
	currentState := *s.state
	s.stateMu.RUnlock()

	// Law I - Create new state (pure function)
	newState := currentState.PlaceNumber(req.Row, req.Col, req.Num)

	// Update atomically
	s.stateMu.Lock()
	s.state = &newState
	s.stateMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"filled":  newState.CountFilled(),
		"valid":   newState.IsValid(),
		"solved":  newState.IsSolved(),
		"version": s.version,
	})
}

// HandleExport exports current board state
func (s *Server) HandleExport(w http.ResponseWriter, r *http.Request) {
	s.stateMu.RLock()
	state := *s.state
	s.stateMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"board":   state.Board,
		"filled":  state.CountFilled(),
		"valid":   state.IsValid(),
		"solved":  state.IsSolved(),
		"version": s.version,
	})
}

// HandleMerge merges incoming board state (Law I associative merge)
func (s *Server) HandleMerge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var incomingState SudokuState
	if err := json.NewDecoder(r.Body).Decode(&incomingState); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("[%s] Merging board with %d filled cells", s.version, incomingState.CountFilled())

	// Law I - Read current state (immutable)
	s.stateMu.RLock()
	currentState := *s.state
	s.stateMu.RUnlock()

	log.Printf("[%s] Current: %d filled, Incoming: %d filled",
		s.version, currentState.CountFilled(), incomingState.CountFilled())

	// Law I - Associative merge
	mergedState := currentState.Merge(incomingState)

	log.Printf("[%s] After merge: %d filled", s.version, mergedState.CountFilled())

	// Update atomically
	s.stateMu.Lock()
	s.state = &mergedState
	s.stateMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"filled":  mergedState.CountFilled(),
		"valid":   mergedState.IsValid(),
		"solved":  mergedState.IsSolved(),
		"message": "Law I: Associative merge completed",
		"version": s.version,
	})
}

// HandleBoard shows current board state (pretty print)
func (s *Server) HandleBoard(w http.ResponseWriter, r *http.Request) {
	s.stateMu.RLock()
	state := *s.state
	s.stateMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"board":   state.Board,
		"filled":  state.CountFilled(),
		"valid":   state.IsValid(),
		"solved":  state.IsSolved(),
		"version": s.version,
	})
}

func (s *Server) Start(addr string) error {
	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// API endpoints
	http.HandleFunc("/place", s.HandlePlace)
	http.HandleFunc("/export", s.HandleExport)
	http.HandleFunc("/merge", s.HandleMerge)
	http.HandleFunc("/board", s.HandleBoard)

	// Web UI endpoints
	http.HandleFunc("/", s.HandleUI)
	http.HandleFunc("/board-html", s.HandleBoardHTML)
	http.HandleFunc("/stats-html", s.HandleStatsHTML)

	log.Printf("[%s] Sudoku server starting on %s", s.version, addr)
	log.Printf("[%s] Law I: Immutable operations (lawtest verified)", s.version)
	log.Printf("[%s] Web UI: http://localhost%s", s.version, addr)
	return http.ListenAndServe(addr, nil)
}

// HandleUI renders the web interface
func (s *Server) HandleUI(w http.ResponseWriter, r *http.Request) {
	s.stateMu.RLock()
	state := *s.state
	s.stateMu.RUnlock()

	stats := web.SudokuStats{
		Filled: state.CountFilled(),
		Valid:  state.IsValid(),
		Solved: state.IsSolved(),
	}

	w.Header().Set("Content-Type", "text/html")
	web.Page(web.SudokuBoard(state.Board), stats, s.version).Render(w)
}

// HandleBoardHTML renders just the board component (for HTMX)
func (s *Server) HandleBoardHTML(w http.ResponseWriter, r *http.Request) {
	s.stateMu.RLock()
	state := *s.state
	s.stateMu.RUnlock()

	w.Header().Set("Content-Type", "text/html")
	web.BoardComponent(web.SudokuBoard(state.Board)).Render(w)
}

// HandleStatsHTML renders just the stats component (for HTMX)
func (s *Server) HandleStatsHTML(w http.ResponseWriter, r *http.Request) {
	s.stateMu.RLock()
	state := *s.state
	s.stateMu.RUnlock()

	stats := web.SudokuStats{
		Filled: state.CountFilled(),
		Valid:  state.IsValid(),
		Solved: state.IsSolved(),
	}

	w.Header().Set("Content-Type", "text/html")
	web.StatsComponent(stats).Render(w)
}
