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

// FindLongest finds the length of longest string in slice.
func FindLongest(ss []string) int {
	l := 0
	for _, s := range ss {
		if i := len(s); i > l {
			l = i
		}
	}
	return l
}
