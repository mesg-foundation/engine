package services

import (
	"encoding/json"

	"github.com/mesg-foundation/core/service"
)

// Get returns a service based on the hash.
func Get(hash string) (service service.Service, err error) {
	db, err := open()
	defer close()
	if err != nil {
		return
	}
	bytes, err := db.Get([]byte(hash), nil)
	if err != nil {
		err = handleErrorNotFound(err, hash)
		return
	}
	err = json.Unmarshal(bytes, &service)
	return
}
