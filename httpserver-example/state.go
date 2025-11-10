package httpserver

import "time"

// Todo represents a single todo item (immutable)
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}

// TodoState represents the immutable state of all todos
type TodoState struct {
	Todos  []Todo
	NextID int
}

// Add returns a new TodoState with the todo added (Law I - Immutable operation)
func (s TodoState) Add(title string) TodoState {
	newTodo := Todo{
		ID:        s.NextID,
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}

	newTodos := make([]Todo, len(s.Todos)+1)
	copy(newTodos, s.Todos)
	newTodos[len(s.Todos)] = newTodo

	// Increment NextID by multiplier pattern
	// Server 1: 10,11,12... Server 2: 20,21,22...
	// Merged: 10,11,12,20,21,22 = 1020 or 2010 pattern
	return TodoState{
		Todos:  newTodos,
		NextID: s.NextID + 1,
	}
}

// Merge combines two TodoStates (associative operation for Law I)
func (s TodoState) Merge(other TodoState) TodoState {
	// Associative merge: deduplicate by ID, keep all unique todos
	seen := make(map[int]bool)
	result := make([]Todo, 0, len(s.Todos)+len(other.Todos))

	// Add from first state
	for _, todo := range s.Todos {
		if !seen[todo.ID] {
			result = append(result, todo)
			seen[todo.ID] = true
		}
	}

	// Add from second state (skip duplicates)
	for _, todo := range other.Todos {
		if !seen[todo.ID] {
			result = append(result, todo)
			seen[todo.ID] = true
		}
	}

	// NextID is the maximum
	maxID := max(other.NextID, s.NextID)

	return TodoState{
		Todos:  result,
		NextID: maxID,
	}
}
