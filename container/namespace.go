package container

// Namespace creates a namespace from a list of string.
func (c *DockerContainer) Namespace(s string) string {
	if s == "" {
		return c.nsprefix
	}
	return c.nsprefix + "-" + s
}
