// Package resolver is a wrapper for system resolver service to call it's tasks.
package resolver

import (
	"errors"
	"fmt"

	"github.com/mesg-foundation/core/api"
)

// Resolvers tasks.
const (
	addPeersTask = "addPeers"
	resolveTask  = "resolve"
)

// Resolver is the system service that responsible from finding the addresses
// of other peers(nodes) in the network that are running the desired services.
// This Resolver is a wrapper for system resolver service to call it's tasks.
type Resolver struct {
	api       *api.API
	serviceID string
}

// New creates a new instance of the Resolver system service.
func New(serviceID string, api *api.API) *Resolver {
	return &Resolver{
		api:       api,
		serviceID: serviceID,
	}
}

// AddPeers attaches peers(nodes) to resolver.
func (r *Resolver) AddPeers(addresses []string) error {
	// TODO: this hack is not something that we should do but
	// it's needed because *parameterValidator is not able to identify
	// string slices for now.
	var addressesInterface []interface{}
	for _, address := range addresses {
		addressesInterface = append(addressesInterface, address)
	}

	// TODO: timeout?
	e, err := r.api.ExecuteAndListen(r.serviceID, addPeersTask, map[string]interface{}{
		"addresses": addressesInterface,
	})
	if err != nil {
		return err
	}

	switch e.Output {
	case "success":
		return nil
	case "error":
		return errors.New(e.OutputData["message"].(string))
	}
	panic("unreachable")
}

// PeerNotFoundError error returned when a peer cannot be found for a service.
type PeerNotFoundError struct {
	ServiceID string
}

func (e *PeerNotFoundError) Error() string {
	return fmt.Sprintf("peer could not found for %q service", e.ServiceID)
}

// Resolve returns the address of a peer(node) that runs the desired service.
func (r *Resolver) Resolve(serviceID string) (peerAddress string, err error) {
	e, err := r.api.ExecuteAndListen(r.serviceID, resolveTask, map[string]interface{}{
		"serviceID": serviceID,
	})
	if err != nil {
		return "", err
	}

	switch e.Output {
	case "found":
		return e.OutputData["address"].(string), nil
	case "notFound":
		return "", &PeerNotFoundError{serviceID}
	case "error":
		return "", errors.New(e.OutputData["message"].(string))
	}
	panic("unreachable")
}
