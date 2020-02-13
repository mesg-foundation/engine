package processsdk

import (
	"fmt"

	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	processpb "github.com/mesg-foundation/engine/process"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/x/process"
)

// SDK is the process sdk.
type SDK struct {
	client *cosmos.Client
}

// New creates a new Process SDK with given options.
func New(client *cosmos.Client) *SDK {
	return &SDK{
		client: client,
	}
}

// Create creates a new process.
func (s *SDK) Create(req *api.CreateProcessRequest) (*processpb.Process, error) {
	acc, err := s.client.GetAccount()
	if err != nil {
		return nil, err
	}
	msg := process.NewMsgCreateProcess(acc.GetAddress(), req)
	tx, err := s.client.BuildAndBroadcastMsg(msg)
	if err != nil {
		return nil, err
	}
	return s.Get(tx.Data)
}

// Delete deletes the process by hash.
func (s *SDK) Delete(req *api.DeleteProcessRequest) error {
	acc, err := s.client.GetAccount()
	if err != nil {
		return err
	}
	msg := process.NewMsgDeleteProcess(acc.GetAddress(), req)
	_, err = s.client.BuildAndBroadcastMsg(msg)
	return err
}

// Get returns the process that matches given hash.
func (s *SDK) Get(hash hash.Hash) (*processpb.Process, error) {
	var p processpb.Process
	if err := s.client.QueryJSON(fmt.Sprintf("custom/%s/%s/%s", process.QuerierRoute, process.QueryGetProcess, hash.String()), nil, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

// List returns all processes.
func (s *SDK) List() ([]*processpb.Process, error) {
	var processes []*processpb.Process
	if err := s.client.QueryJSON(fmt.Sprintf("custom/%s/%s", process.QuerierRoute, process.QueryListProcesses), nil, &processes); err != nil {
		return nil, err
	}
	return processes, nil
}
