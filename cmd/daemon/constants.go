package daemon

const (
	name          = "mesg-daemon"
	sharedNetwork = "mesg-shared-network"
	image         = "mesg/daemon"
	socketPath    = "/var/run/docker.sock"
)
