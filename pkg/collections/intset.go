package collections

// IntSet is a map that replicates a set for int32s from languages with generics
type IntSet map[int]struct{}

// Add a value to the set
func (s IntSet) Add(key int) {
	s[key] = void
}

// Contains checks if the key belongs in the map
func (s IntSet) Contains(key int) bool {
	_, exists := s[key]
	return exists
}

// Each executes a function for every item in the set
func (s IntSet) Each(f func(key int)) {
	for key := range s {
		f(key)
	}
}
