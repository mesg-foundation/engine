package services

import (
	"github.com/mesg-foundation/core/service"
)

// All returns all deployed services.
func All() ([]*service.Service, error) {
	db, err := open()
	defer close()
	if err != nil {
		return nil, err
	}
	var services []*service.Service
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
			return nil, err
		service, err := decode(string(iter.Key()), iter.Value())
		if err != nil {
		}
		services = append(services, &service)
	}
	iter.Release()
	return services, iter.Error()
}
