package services

import (
	"encoding/json"

	"github.com/mesg-foundation/core/service"
)

// encode returns the marshaled version of the service
func encode(service *service.Service) ([]byte, error) {
	return json.Marshal(service)
}

// decode decodes the data and return an DecodeError if not possible
func decode(hash string, data []byte) (*service.Service, error) {
	s := &service.Service{}
	if err := json.Unmarshal(data, s); err != nil {
		return nil, &DecodeError{Hash: hash}
	}
	return s, nil
}
