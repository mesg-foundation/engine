package importer

import (
	"github.com/mesg-foundation/core/service"

	yaml "gopkg.in/yaml.v2"
)

// From imports a service from a source
func From(source string) (importedService *service.Service, err error) {
	return fromPath(source)
}

// fromPath imports a service from a path
func fromPath(path string) (importedService *service.Service, err error) {
	isValid, err := IsValid(path)
	if err != nil {
		return
	}
	if isValid == false {
		err = &ValidationError{}
		return
	}
	data, err := readServiceFile(path)
	if err != nil {
		return
	}
	importedService = &service.Service{}
	err = yaml.UnmarshalStrict(data, importedService)
	return
}
