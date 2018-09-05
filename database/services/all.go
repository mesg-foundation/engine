package services

import (
	"encoding/json"

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
		var service service.Service
		if err := json.Unmarshal(iter.Value(), &service); err != nil {
			return nil, err
		}
		services = append(services, &service)
	}
	iter.Release()
	return services, iter.Error()
}
