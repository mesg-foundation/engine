package services

import (
	"github.com/golang/protobuf/proto"
	"github.com/mesg-foundation/core/service"
)

// All returns all deployed services
func All() (services []*service.Service, err error) {
	db, err := open()
	defer close()
	if err != nil {
		return
	}
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		var service service.Service
		err = proto.Unmarshal(iter.Value(), &service)
		if err != nil {
			return
		}
		services = append(services, &service)
	}
	iter.Release()
	err = iter.Error()
	return
}
