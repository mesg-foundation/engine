package services

import (
	"github.com/mesg-foundation/core/service"
	"github.com/sirupsen/logrus"
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
		service, err := decode(string(iter.Key()), iter.Value())
		if err != nil {
			// Ignore all decode errors (possibly due to a service structure change or database corruption)
			if decodeErr, ok := err.(*DecodeError); ok {
				logrus.WithField("service", decodeErr.Hash).Warning(decodeErr.Error())
			} else {
				return nil, err
			}
		} else {
			services = append(services, service)
		}
	}
	iter.Release()
	return services, iter.Error()
}
