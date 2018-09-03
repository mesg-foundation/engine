package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"
)

// ColorizeJSON adds colors to a stringified JSON. If an error happens, the original string is returned.
func ColorizeJSON(data string) string {
	var decoded map[string]interface{}
	if err := json.Unmarshal([]byte(data), &decoded); err != nil {
		return data
	}
	var outputs []string
	for key, value := range decoded {
		outputs = append(outputs, fmt.Sprintf("%v = %v", aurora.Cyan(key), value))
	}
	return strings.Join(outputs, ", ")
}
