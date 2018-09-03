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

// SliceIndex returns the index e in a.
func SliceIndex(a []string, e string) int {
	for i, s := range a {
		if s == e {
			return i
		}
	}
	return 0
}
