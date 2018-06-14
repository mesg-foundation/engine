package services

import (
	"github.com/golang/protobuf/proto"
	"github.com/mesg-foundation/core/service"
)

// Get a service based on it's hash
func Get(hash string) (service service.Service, err error) {
	db, err := open()
	defer close()
	if err != nil {
		return
	}
	bytes, err := db.Get([]byte(hash), nil)
	err = handleErrorNotFound(err, hash)
	if err != nil {
		return
	}
	err = proto.Unmarshal(bytes, &service)
	return
}
