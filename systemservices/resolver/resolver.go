package resolver

import "github.com/mesg-foundation/core/api"

// Resolver is the system service responsible for getting the address of
// other node of the network that are running the desired service.
type Resolver struct{}

// New creates a new instance of the Resolver system service.
// It communicates with the service by only using the API package.
// serviceID of the associated MESG Service.
// api is an instance to the api package required to communication with the service.
func New(serviceID string, api *api.API) (*Resolver, error) {
	return nil, nil
}
