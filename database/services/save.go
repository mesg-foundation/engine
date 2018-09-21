package services

import (
	"github.com/mesg-foundation/core/service"
)

// Save stores a service in the database.
func Save(service *service.Service) error {
	bytes, err := encode(service)
	if err != nil {
		return err
	}
	db, err := open()
	defer close()
	if err != nil {
		return err
	}
	return db.Put([]byte(service.ID), bytes, nil)
}
