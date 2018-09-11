package xstrings

// SliceContains returns true if slice a contains e element, false otherwise.
func SliceContains(a []string, e string) bool {
	for _, s := range a {
		if s == e {
			return true
		}
	}
	return false
}

// AppendSpace appends n times space to s.
func AppendSpace(s string, n int) string {
	if n <= 0 {
		return s
	}
	for i := 0; i < n; i++ {
		s += " "
	}
	return s
}
