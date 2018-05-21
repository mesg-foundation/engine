package dependency

type dependency interface {
	GetPorts() []string
	GetVolumes() []string
	GetVolumesfrom() []string
}
type service interface {
	Namespace() string
	GetDependencies() map[string]dependency
}

// StatusType of the service
type StatusType uint

// status for services
const (
	STOPPED StatusType = 0
	RUNNING StatusType = 1
	PARTIAL StatusType = 2
)
