package example

type Item struct {
	ID   string
	Name string
}

type Items []Item

// Merge combines two slices
func (items Items) Merge(other Items) Items {
	result := make(Items, len(items)+len(other))
	copy(result, items)
	copy(result[len(items):], other)
	return result
}

// Append adds an item
func (items Items) Append(item Item) Items {
	return append(items, item)
}
