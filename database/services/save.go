package services

import (
	"github.com/golang/protobuf/proto"
	"github.com/mesg-foundation/core/service"
)

// Save a service in the database and return the hash or the error if something wrong happened
func Save(service *service.Service) (hash string, err error) {
	bytes, err := proto.Marshal(service)
	if err != nil {
		return "", err
	}
	db, err := open()
	defer close()
	if err != nil {
		return "", err
	}
	return service.Hash(), db.Put([]byte(hash), bytes, nil)
}
