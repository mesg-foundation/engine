package instance

import "github.com/mesg-foundation/engine/hash"

// Instance is an instance of a service.
// This contains a reference to a service that is running.
// Multiple instances can run for the same service as long as they have different configurations
type Instance struct {
	Hash        hash.Hash
	ServiceHash hash.Hash
}
