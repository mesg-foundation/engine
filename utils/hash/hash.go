package hash

import (
	"strings"
)

const separator = "."

// Calculate will return a hash according to the data given
func Calculate(data []string) string {
	return strings.Join(data, separator)
}
