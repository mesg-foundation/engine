package serialize

import (
	s "github.com/mesg-foundation/core/service"

	yaml "gopkg.in/yaml.v2"
)

// FromPath imports a service from a path
func FromPath(path string) (service *s.Service, err error) {
	isValid, err := IsValidFromPath(path)
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
	service = &s.Service{}
	err = yaml.UnmarshalStrict(data, service)
	return
}
