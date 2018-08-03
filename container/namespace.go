package container

import (
	"strings"
)

const (
	namespacePrefix    string = "mesg"
	namespaceSeparator string = "-"
)

// Namespace creates a namespace from a list of string.
func Namespace(ss []string) string {
	ssWithPrefix := append([]string{namespacePrefix}, ss...)
	namespace := strings.Join(ssWithPrefix, namespaceSeparator)
	namespace = strings.Replace(namespace, " ", namespaceSeparator, -1)
	return namespace
}
