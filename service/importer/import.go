package importer

import (
	yaml "gopkg.in/yaml.v2"
)

// From imports a service from a source.
func From(source string) (*ServiceDefinition, error) {
	return fromPath(source)
}

// fromPath imports a service from a path.
func fromPath(path string) (*ServiceDefinition, error) {
	isValid, err := IsValid(path)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, &ValidationError{}
	}
	data, err := readServiceFile(path)
	if err != nil {
		return nil, err
	}
	var importedService ServiceDefinition
	return &importedService, yaml.UnmarshalStrict(data, &importedService)
}
