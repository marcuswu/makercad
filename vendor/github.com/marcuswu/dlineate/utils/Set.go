package utils

var exists = struct{}{}

// Set a set of unsigned integers
type Set struct {
	m map[uint]struct{}
}

// NewSet create a new set
func NewSet() *Set {
	s := &Set{}
	s.m = make(map[uint]struct{})
	return s
}

// Add Adds a value to the set
func (s *Set) Add(value uint) {
	s.m[value] = exists
}

// Remove Removes a value from the set
func (s *Set) Remove(value uint) {
	delete(s.m, value)
}

// Contains Return whether a set contains a value
func (s *Set) Contains(value uint) bool {
	_, c := s.m[value]
	return c
}

// AddList Adds a list of values to the set
func (s *Set) AddList(values []uint) {
	for _, v := range values {
		s.Add(v)
	}
}

// AddSet Adds the contents of another set to this set
func (s *Set) AddSet(values *Set) {
	s.AddList(values.Contents())
}

// Contents returns a copy of the underlying set data
func (s *Set) Contents() []uint {
	keys := make([]uint, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}
	return keys
}

// Count returns the number of elements in the set
func (s *Set) Count() int {
	return len(s.m)
}

func (s *Set) Clear() {
	for value := range s.m {
		delete(s.m, value)
	}
}
