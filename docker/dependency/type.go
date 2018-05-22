package dependency

type Dependency interface {
	GetPorts() []string
	GetVolumes() []string
	GetVolumesfrom() []string
}
type Service interface {
	Namespace() string
	GetDependency(name string) Dependency
}
