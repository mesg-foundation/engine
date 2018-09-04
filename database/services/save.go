package services

import (
	"encoding/json"

	"github.com/mesg-foundation/core/service"
)

// Save stores a service in the database and returns a hash or an error.
func Save(service *service.Service) error {
	bytes, err := json.Marshal(service)
	if err != nil {
		return err
	}
	db, err := open()
	defer close()
	if err != nil {
		return err
	}

	// TODO(ilgooz) rm this when we have a New() initializer for Service type.
	service.ID = service.Hash()

	return db.Put([]byte(service.ID), bytes, nil)
}
