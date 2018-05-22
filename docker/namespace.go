package docker

import (
	"strings"
)

const namespacePrefix string = "MESG"
const namespaceSeparator string = "-"

// Namespace creates a namespace from a list of string
func Namespace(ss []string) string {
	names := append([]string{namespacePrefix}, ss...)
	name := strings.Join(names, namespaceSeparator)
	name = strings.Replace(name, " ", namespaceSeparator, -1)
	return name
}
