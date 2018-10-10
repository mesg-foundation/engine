package resolver

import (
	"errors"
	"fmt"

	"github.com/mesg-foundation/core/api"
)

// Resolvers tasks.
const (
	addPeersTask string = "AddPeers"
	resolveTask         = "Resolve"
)

// Resolver is the system service responsible for getting the address of
// other node of the network that are running the desired service.
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

// AddPeers is the task that actually new add peers to the service.
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
		return errors.New(e.OutputData["error"].(string))
	default:
		return fmt.Errorf("unexpected output %s", e.Output)
	}
}

// ResolveFoundOutput is the found output data of resolve task.
type ResolveFoundOutput struct {
	Address string `json:"address"` // Address is the IP address of core peer.
}

// Resolve is the task that return the address of a core that runs the desired service.
func (r *Resolver) Resolve(serviceID string) (address string, err error) {
	e, err := r.api.ExecuteAndListen(r.serviceID, resolveTask, map[string]interface{}{
		"serviceID": serviceID,
	})

	switch e.Output {
	case "found":
		return e.OutputData["found"].(*ResolveFoundOutput).Address, nil
	case "notFound":
		return "", fmt.Errorf("address for service id %s not found", serviceID)
	case "error":
		return "", errors.New(e.OutputData["error"].(string))
	default:
		return "", fmt.Errorf("unexpected output %s", e.Output)
	}
}
