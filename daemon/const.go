package daemon

const (
	name         = "core"
	dockerSocket = "/var/run/docker.sock"
	volume       = "mesg-core"
	path         = "/mesg"
)

// Namespace returns the namespace of the MESG Core.
func Namespace() []string {
	return []string{name}
}
