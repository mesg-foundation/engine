package ownershipsdk

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/ownership"
)

// SDK is the ownership sdk.
type SDK struct {
	cdc    *codec.Codec
	client *cosmos.Client
}

// New returns the ownership sdk.
func New(cdc *codec.Codec, client *cosmos.Client) *SDK {
	sdk := &SDK{
		cdc:    cdc,
		client: client,
	}
	return sdk
}

// List returns all ownerships.
func (s *SDK) List() ([]*ownership.Ownership, error) {
	var ownerships []*ownership.Ownership
	if err := s.client.Query("custom/"+backendName+"/list", &ownerships); err != nil {
		return nil, err
	}
	return ownerships, nil
}
