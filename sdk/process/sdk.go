package processsdk

import (
	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/api"
)

// SDK is the process sdk.
type SDK struct {
	client     *cosmos.Client
	accAddress cosmostypes.AccAddress
}

// New creates a new Process SDK with given options.
func New(client *cosmos.Client, accAddress cosmostypes.AccAddress) *SDK {
	return &SDK{
		client:     client,
		accAddress: accAddress,
	}
}

// Create creates a new process.
func (s *SDK) Create(req *api.CreateProcessRequest) (*process.Process, error) {
	msg := newMsgCreateProcess(s.accAddress, req)
	tx, err := s.client.BuildAndBroadcastMsg(msg)
	if err != nil {
		return nil, err
	}
	return s.Get(tx.Data)
}

// Delete deletes the process by hash.
func (s *SDK) Delete(req *api.DeleteProcessRequest) error {
	msg := newMsgDeleteProcess(s.accAddress, req)
	_, err := s.client.BuildAndBroadcastMsg(msg)
	return err
}

// Get returns the process that matches given hash.
func (s *SDK) Get(hash hash.Hash) (*process.Process, error) {
	var process process.Process
	if err := s.client.Query("custom/"+backendName+"/get/"+hash.String(), nil, &process); err != nil {
		return nil, err
	}
	return &process, nil
}

// List returns all processes.
func (s *SDK) List() ([]*process.Process, error) {
	var processes []*process.Process
	if err := s.client.Query("custom/"+backendName+"/list", nil, &processes); err != nil {
		return nil, err
	}
	return processes, nil
}
