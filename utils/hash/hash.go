package hash

import (
	"strings"
)

const separator = "."

// Calculate will return a hash according to the data given
func Calculate(data []string) (res string) {
	res = strings.Join(data, separator)
	return
}
