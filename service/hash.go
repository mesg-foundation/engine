package service

import "github.com/mesg-foundation/core/x/xstructhash"

// Hash calculates and returns the hash of the service.
func (service *Service) Hash() string {
	return xstructhash.Hash(service, 1)
}
