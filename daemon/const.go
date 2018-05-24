package daemon

const (
	name          = "daemon"
	sharedNetwork = "shared"
	image         = "mesg/daemon"
	dockerSocket  = "/var/run/docker.sock"
)

// Namespace returns the namespace of the daemon
func Namespace() []string {
	return []string{name}
}

// NamespaceNetwork returns the namespace of the daemon shared network
func NamespaceNetwork() []string {
	return []string{sharedNetwork}
}
