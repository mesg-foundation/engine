package resultsdk

import (
	"context"

	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/result"
)

// SDK is the result sdk.
type SDK struct {
	client *cosmos.Client
	kb     *cosmos.Keybase
}

// New returns the result sdk.
func New(client *cosmos.Client, kb *cosmos.Keybase) *SDK {
	sdk := &SDK{
		client: client,
		kb:     kb,
	}
	return sdk
}

// Create creates a new result.
func (s *SDK) Create(req *api.CreateResultRequest, accountName, accountPassword string) (*result.Result, error) {
	acc, err := s.kb.Get(accountName)
	if err != nil {
		return nil, err
	}
	msg := newMsgCreateResult(req, acc.GetAddress())
	tx, err := s.client.BuildAndBroadcastMsg(msg, accountName, accountPassword)
	if err != nil {
		return nil, err
	}
	return s.Get(tx.Data)
}

// Get returns the result that matches given hash.
func (s *SDK) Get(hash hash.Hash) (*result.Result, error) {
	var res *result.Result
	if err := s.client.Query("custom/"+backendName+"/get/"+hash.String(), nil, &res); err != nil {
		return nil, err
	}
	return res, nil
}

// List returns all results.
func (s *SDK) List() ([]*result.Result, error) {
	var ress []*result.Result
	if err := s.client.Query("custom/"+backendName+"/list", nil, &ress); err != nil {
		return nil, err
	}
	return ress, nil
}

// Stream returns result that matches given hash.
func (s *SDK) Stream(ctx context.Context, req *api.StreamResultRequest) (chan *result.Result, chan error, error) {
	// TODO:
	// if err := req.Filter.Validate(); err != nil {
	// 	return nil, nil, err
	// }
	stream, serrC, err := s.client.Stream(ctx, cosmos.EventModuleQuery(backendName))
	if err != nil {
		return nil, nil, err
	}
	execC := make(chan *result.Result)
	errC := make(chan error)
	go func() {
	loop:
		for {
			select {
			case hash := <-stream:
				exec, err := s.Get(hash)
				if err != nil {
					errC <- err
				}
				// TODO:
				// if req.Filter.Match(exec) {
				execC <- exec
				// }
			case err := <-serrC:
				errC <- err
			case <-ctx.Done():
				break loop
			}
		}
		close(errC)
		close(execC)
	}()
	return execC, errC, nil
}
