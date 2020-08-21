package set

// StringSet implements minimal set-like behaviour
type StringSet map[string]bool

// Contains returns true if `item` is a member of the set
func (s StringSet) Contains(item string) bool {
	_, ok := s[item]
	return ok
}

// FromSlice creates a new StringSet from the specified slice
func FromSlice(items []string) StringSet {
	res := make(StringSet, len(items))
	for _, item := range items {
		res[item] = true
	}
	return res
}
