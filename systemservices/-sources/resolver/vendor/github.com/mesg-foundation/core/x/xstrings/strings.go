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

// SliceIndex returns the index e in a, return -1 if not found.
func SliceIndex(a []string, e string) int {
	for i, s := range a {
		if s == e {
			return i
		}
	}
	return -1
}
