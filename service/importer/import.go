package importer

import (
	"github.com/mesg-foundation/core/service"
	yaml "gopkg.in/yaml.v2"
)

// From imports a service from a source
func From(source string) (*service.Service, error) {
	return fromPath(source)
}

// fromPath imports a service from a path
func fromPath(path string) (*service.Service, error) {
	isValid, err := IsValid(path)
	if err != nil {
		return nil, err
	}
	if isValid == false {
		return nil, &ValidationError{}
	}
	data, err := readServiceFile(path)
	if err != nil {
		return nil, err
	}
	var importedService service.Service
	if err := yaml.UnmarshalStrict(data, importedService); err != nil {
		return nil, err
	}
	return &importedService, nil
}
