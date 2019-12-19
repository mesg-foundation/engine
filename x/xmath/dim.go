package xmath

// MaxUint64 returns the larger of a and b...
func MaxUint64(a uint64, b ...uint64) uint64 {
	for _, i := range b {
		if i > a {
			a = i
		}
	}
	return a
}
