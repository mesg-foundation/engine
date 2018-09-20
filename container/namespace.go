package container

import (
	"strings"
)

const namespaceSeparator string = "-"

// Namespace creates a namespace from a list of string.
func (c *Container) Namespace(ss []string) string {
	ssWithPrefix := append([]string{c.config.Core.Name}, ss...)
	namespace := strings.Join(ssWithPrefix, namespaceSeparator)
	namespace = strings.Replace(namespace, " ", namespaceSeparator, -1)
	return namespace
}
