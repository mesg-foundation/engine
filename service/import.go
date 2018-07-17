package service

import (
	"errors"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

func readFromPath(path string) (data []byte, err error) {
	file := filepath.Join(path, "mesg.yml")
	data, err = ioutil.ReadFile(file)
	return
}

// ImportFromPath returns the service of the given path
func ImportFromPath(path string) (service *Service, err error) {
	validation, err := ValidateService(path)
	if err != nil {
		return
	}
	if validation.IsValid() == false {
		err = errors.New("Service is not valid. Run the command 'service validate' for more details")
		return
	}
	data, err := readFromPath(path)
	if err != nil {
		return
	}
	service = &Service{}
	err = yaml.UnmarshalStrict(data, service)
	return
}
