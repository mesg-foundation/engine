package service

import (
	"github.com/cnf/structhash"
)

// Hash calculate and return the hash of the service
func (service *Service) Hash() (hash string, err error) {
	hash, err = structhash.Hash(service, 1)
	return
}
