package services

import (
	"github.com/golang/protobuf/proto"
	"github.com/mesg-foundation/core/service"
)

// Save stores a service in the database and returns a hash or an error.
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
	hash = service.Hash()
	return hash, db.Put([]byte(hash), bytes, nil)
}
