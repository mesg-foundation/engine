package services

import (
	"github.com/golang/protobuf/proto"
	"github.com/mesg-foundation/core/service"
)

// Save a service in the database and return the hash or the error if something wrong happened
func Save(service *service.Service) (hash string, err error) {
	bytes, err := proto.Marshal(service)
	if err != nil {
		return
	}
	hash, err = service.Hash()
	if err != nil {
		return
	}
	db := open()
	defer close()
	err = db.Put([]byte(hash), bytes, nil)
	return
}
