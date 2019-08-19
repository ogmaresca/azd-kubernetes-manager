package collections

var void struct{}

// StringSet is a map that replicates a set for strings from languages with generics
type StringSet map[string]struct{}

// Add a value to the set
func (s StringSet) Add(key string) {
	s[key] = void
}

// Contains checks if the key belongs in the map
func (s StringSet) Contains(key string) bool {
	_, exists := s[key]
	return exists
}

// Each executes a function for every item in the set
func (s StringSet) Each(f func(key string)) {
	for key := range s {
		f(key)
	}
}
