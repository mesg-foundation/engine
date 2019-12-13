package ownershipsdk

import (
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/ownership"
)

// SDK is the ownership sdk.
type SDK struct {
	client *cosmos.Client
}

// New returns the ownership sdk.
func New(client *cosmos.Client) *SDK {
	sdk := &SDK{
		client: client,
	}
	return sdk
}

// List returns all ownerships.
func (s *SDK) List() ([]*ownership.Ownership, error) {
	var ownerships []*ownership.Ownership
	if err := s.client.Query("custom/"+backendName+"/list", nil, &ownerships); err != nil {
		return nil, err
	}
	return ownerships, nil
}
