package instance

import "github.com/mesg-foundation/engine/hash"

func New(serviceHash hash.Hash, envHash hash.Hash) *Instance {
	inst := &Instance{
		ServiceHash: serviceHash,
		EnvHash:     envHash,
	}
	inst.Hash = hash.Dump(inst)
	return inst
}
