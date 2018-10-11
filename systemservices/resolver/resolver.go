package resolver

import (
	"errors"
	"fmt"

	"github.com/mesg-foundation/core/api"
)

// Resolvers tasks.
const (
	addPeersTask = "AddPeers"
	resolveTask  = "Resolve"
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
	e, err := r.api.ExecuteAndListen(r.serviceID, addPeersTask, map[string]interface{}{
		"addresses": addresses,
	})
	if err != nil {
		return err
	}

	switch e.Output {
	case "success":
		return nil
	case "error":
		return errors.New(e.OutputData["message"].(string))
	default:
		return fmt.Errorf("unexpected output %s", e.Output)
	}
}

// Resolve return the address of a peer(node) that runs the desired service.
func (r *Resolver) Resolve(serviceID string) (address string, err error) {
	e, err := r.api.ExecuteAndListen(r.serviceID, resolveTask, map[string]interface{}{
		"serviceID": serviceID,
	})

	switch e.Output {
	case "found":
		return e.OutputData["address"].(string), nil
	case "notFound":
		return "", fmt.Errorf("address for service id %s not found", serviceID)
	case "error":
		return "", errors.New(e.OutputData["message"].(string))
	default:
		return "", fmt.Errorf("unexpected output %s", e.Output)
	}
}
