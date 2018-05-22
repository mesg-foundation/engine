package service

type service interface {
	Namespace() string
	GetDependenciesKeys() []string
}
