package services

import (
	"github.com/golang/protobuf/proto"
	"github.com/mesg-foundation/core/service"
)

// Get a service based on it's hash
func Get(hash string) (service service.Service, err error) {
	db := open()
	defer close()
	bytes, err := db.Get([]byte(hash), nil)
	if err != nil {
		return
	}
	err = proto.Unmarshal(bytes, &service)
	return
}
