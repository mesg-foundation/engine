package instance

import (
	"github.com/mesg-foundation/engine/ext/xvalidator"
	"github.com/mesg-foundation/engine/hash"
)

// New returns a new Instance and validate it.
func New(serviceHash hash.Hash, envHash hash.Hash) (*Instance, error) {
	inst := &Instance{
		ServiceHash: serviceHash,
		EnvHash:     envHash,
	}
	inst.Hash = hash.Dump(inst)
	return inst, xvalidator.Struct(inst)
}
