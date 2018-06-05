package service

import (
	"github.com/cnf/structhash"
)

// Hash calculate and return the hash of the service
func (service *Service) Hash() (hash string) {
	// Ignore the err result because the lib always return nil
	hash, _ = structhash.Hash(service, 1)
	return
}
