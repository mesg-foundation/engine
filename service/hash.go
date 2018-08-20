package service

import (
	"github.com/cnf/structhash"
)

// Hash calculates and returns the hash of the service.
func (service *Service) Hash() string {
	// Ignore the err result because the lib always return nil
	hash, _ := structhash.Hash(service, 1) // TODO: why not reuse the package utils/hash?
	return hash
}
