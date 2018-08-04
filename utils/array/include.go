package array

// IncludedIn returns true if the element value is in the string slice, false otherwise.
func IncludedIn(arr []string, value string) bool {
	if len(arr) == 0 {
		return false
	}
	i := 0
	for _, item := range arr {
		if item == value {
			break
		}
		i++
	}
	return i != len(arr)
}
