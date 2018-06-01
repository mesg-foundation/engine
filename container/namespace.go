package container

import (
	"strings"
)

const (
	namespacePrefix    string = "mesg"
	namespaceSeparator string = "-"
	serviceTagPrefix   string = "mesg/"
)

// Namespace creates a namespace from a list of string
func Namespace(ss []string) string {
	ssWithPrefix := append([]string{namespacePrefix}, ss...)
	namespace := strings.Join(ssWithPrefix, namespaceSeparator)
	namespace = strings.Replace(namespace, " ", namespaceSeparator, -1)
	return namespace
}

// ServiceTag returns the tag for a docker image for a service
func ServiceTag(ss []string) string {
	namespace := strings.Join(ss, namespaceSeparator)
	namespace = strings.Replace(namespace, " ", namespaceSeparator, -1)
	return serviceTagPrefix + namespace
}
