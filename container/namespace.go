package container

import "github.com/mesg-foundation/engine/version"

// Namespace creates a namespace from a list of string.
func (c *DockerContainer) Namespace(s string) string {
	if s == "" {
		return version.Name
	}
	return version.Name + "-" + s
}
