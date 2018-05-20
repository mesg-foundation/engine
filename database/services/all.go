package services

import (
	"github.com/golang/protobuf/proto"
	"github.com/mesg-foundation/core/service"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// All returns all deployed services
func All() (services []*service.Service, err error) {
	db := open()
	defer close()
	iter := db.NewIterator(util.BytesPrefix([]byte("")), nil)
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
