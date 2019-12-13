package instancesdk

import (
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/instance"
	"github.com/mesg-foundation/engine/protobuf/api"
)

// SDK is the instance sdk.
type SDK struct {
	client *cosmos.Client
}

// New returns the instance sdk.
func New(client *cosmos.Client) *SDK {
	sdk := &SDK{
		client: client,
	}
	return sdk
}

// Get returns the instance that matches given hash.
func (s *SDK) Get(hash hash.Hash) (*instance.Instance, error) {
	var instance instance.Instance
	if err := s.client.Query("custom/"+backendName+"/get/"+hash.String(), nil, &instance); err != nil {
		return nil, err
	}
	return &instance, nil
}

// List returns all instances.
func (s *SDK) List(f *api.ListInstanceRequest_Filter) ([]*instance.Instance, error) {
	var instances []*instance.Instance
	if err := s.client.Query("custom/"+backendName+"/list", f, &instances); err != nil {
		return nil, err
	}
	return instances, nil
}

// Exists returns if a instance already exists.
func (s *SDK) Exists(hash hash.Hash) (bool, error) {
	var exists bool
	if err := s.client.Query("custom/"+backendName+"/exists/"+hash.String(), nil, &exists); err != nil {
		return false, err
	}
	return exists, nil
}
