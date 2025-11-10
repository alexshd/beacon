package example

// Add combines two integers
func Add(a, b int) int {
	return a + b
}

// Concat combines two strings
func Concat(a, b string) string {
	return a + b
}

// MergeMap combines two maps (not comparable)
func MergeMap(a, b map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range a {
		result[k] = v
	}
	for k, v := range b {
		result[k] = v
	}
	return result
}

type State struct {
	count int
}

// Merge combines two states
func (s State) Merge(other State) State {
	return State{count: s.count + other.count}
}
