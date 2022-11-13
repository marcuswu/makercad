package utils

// StringSet a set of unsigned integers
type StringSet struct {
	m map[string]struct{}
}

// NewStringSet create a new set
func NewStringSet() *StringSet {
	s := &StringSet{}
	s.m = make(map[string]struct{})
	return s
}

// Add Adds a value to the set
func (s *StringSet) Add(value string) {
	s.m[value] = exists
}

// Remove Removes a value from the set
func (s *StringSet) Remove(value string) {
	delete(s.m, value)
}

// Contains Return whether a set contains a value
func (s *StringSet) Contains(value string) bool {
	_, c := s.m[value]
	return c
}

// AddList Adds a list of values to the set
func (s *StringSet) AddList(values []string) {
	for _, v := range values {
		s.Add(v)
	}
}

// AddStringSet Adds the contents of another set to this set
func (s *StringSet) AddStringSet(values *StringSet) {
	s.AddList(values.Contents())
}

// Contents returns a copy of the underlying set data
func (s *StringSet) Contents() []string {
	keys := make([]string, 0, len(s.m))
	for k := range s.m {
		keys = append(keys, k)
	}
	return keys
}

// Count returns the number of elements in the set
func (s *StringSet) Count() int {
	return len(s.m)
}

func (s *StringSet) Clear() {
	for value := range s.m {
		delete(s.m, value)
	}
}
