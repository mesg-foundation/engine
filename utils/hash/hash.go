package hash

import (
	"strings"
)

const separator = "."

// Calculate returns a hash according to the given data.
func Calculate(data []string) string {
	return strings.Join(data, separator)
}
