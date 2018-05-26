package daemon

const (
	name          = "daemon"
	sharedNetwork = "shared"
	networkDriver = "overlay"
	dockerSocket  = "/var/run/docker.sock"
)

// Namespace returns the namespace of the daemon
func namespace() []string {
	return []string{name}
}

// NamespaceNetwork returns the namespace of the daemon shared network
func namespaceNetwork() []string {
	return []string{sharedNetwork}
}
