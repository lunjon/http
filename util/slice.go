package util

// Map slice a using f.
func Map(a []string, f func(string) string) []string {
	if len(a) == 0 {
		return []string{}
	}

	s := make([]string, len(a))
	for i, v := range a {
		s[i] = f(v)
	}
	return s
}

// Filter returns all values of a for which f returns true.
func Filter(a []string, f func(string) bool) []string {
	if len(a) == 0 {
		return []string{}
	}

	s := make([]string, 0)
	for _, v := range a {
		if f(v) {
			s = append(s, v)
		}
	}
	return s
}
