package services

import (
	"encoding/json"

	"github.com/mesg-foundation/core/service"
)

// Get returns a service based on the id.
func Get(id string) (*service.Service, error) {
	db, err := open()
	defer close()
	if err != nil {
		return nil, err
	}
	bytes, err := db.Get([]byte(id), nil)
	if err != nil {
		err = handleErrorNotFound(err, id)
		return nil, err
	}
	s := &service.Service{}
	return s, json.Unmarshal(bytes, s)
}
