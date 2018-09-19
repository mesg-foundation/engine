package container

import (
	"log"
	"strings"

	"github.com/mesg-foundation/core/config"
)

const namespaceSeparator string = "-"

// Namespace creates a namespace from a list of string.
// TODO: Put this function as in the Container type.
// TODO: load config in the container struct
// TODO: change input as "..string"
func Namespace(ss []string) string {
	c, err := config.Global()
	if err != nil {
		log.Fatalln("TODO: Manage error")
	}
	ssWithPrefix := append([]string{c.Core.Name}, ss...)
	namespace := strings.Join(ssWithPrefix, namespaceSeparator)
	namespace = strings.Replace(namespace, " ", namespaceSeparator, -1)
	return namespace
}
