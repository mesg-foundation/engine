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
		var svc service.Service
		err = proto.Unmarshal(iter.Value(), &svc)
		if err != nil {
			return
		}
		services = append(services, &svc)
	}
	iter.Release()
	err = iter.Error()
	return
}
