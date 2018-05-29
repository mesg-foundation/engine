package daemon

const (
	name         = "daemon"
	dockerSocket = "/var/run/docker.sock"
)

// Namespace returns the namespace of the daemon
func Namespace() []string {
	return []string{name}
}
