package service

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// ImportFromFile will return a service associated to the file given in parameter
func ImportFromFile(filename string) (service *Service, err error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	service = &Service{}
	err = yaml.Unmarshal(file, service)
	if err != nil {
		return
	}
	return
}
