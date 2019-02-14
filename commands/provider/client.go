package provider

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mesg-foundation/core/protobuf/acknowledgement"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	uuid "github.com/satori/go.uuid"
)

// client is a wrapper for core client that
// sets up stream connection and provides more frienldy api.
type client struct {
	client coreapi.CoreClient
}

// ListenEvents returns a channel with event data streaming.
func (c *client) ListenEvents(id, eventFilter string) (chan *coreapi.EventData, chan error, error) {
	stream, err := c.client.ListenEvent(context.Background(), &coreapi.ListenEventRequest{
		ServiceID:   id,
		EventFilter: eventFilter,
	})
	if err != nil {
		return nil, nil, err
	}

	resultC := make(chan *coreapi.EventData)
	errC := make(chan error)

	go func() {
		<-stream.Context().Done()
		errC <- stream.Context().Err()
	}()
	go func() {
		for {
			if res, err := stream.Recv(); err != nil {
				errC <- err
				break
			} else {
				resultC <- res
			}
		}
	}()

	if err := acknowledgement.WaitForStreamToBeReady(stream); err != nil {
		return nil, nil, err
	}

	return resultC, errC, nil
}

// ListenResults returns a channel with event results streaming..
func (c *client) ListenResult(id, taskFilter, outputFilter string, tagFilters []string) (chan *coreapi.ResultData, chan error, error) {
	resultC := make(chan *coreapi.ResultData)
	errC := make(chan error)

	stream, err := c.client.ListenResult(context.Background(), &coreapi.ListenResultRequest{
		ServiceID:    id,
		TaskFilter:   taskFilter,
		OutputFilter: outputFilter,
		TagFilters:   tagFilters,
	})
	if err != nil {
		return nil, nil, err
	}

	go func() {
		<-stream.Context().Done()
		errC <- stream.Context().Err()
	}()
	go func() {
		for {
			if res, err := stream.Recv(); err != nil {
				errC <- err
				break
			} else {
				resultC <- res
			}
		}
	}()

	if err := acknowledgement.WaitForStreamToBeReady(stream); err != nil {
		return nil, nil, err
	}

	return resultC, errC, nil
}

// ExecuteTask executes task on given service.
func (c *client) ExecuteTask(id, taskKey, inputData string, tags []string) error {
	_, err := c.client.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID:     id,
		TaskKey:       taskKey,
		InputData:     inputData,
		ExecutionTags: tags,
	})
	return err
}

// ExecuteAndListen executes task and listens for it's results.
func (c *client) ExecuteAndListen(id, taskKey string, inputData interface{}) (*coreapi.ResultData, error) {
	data, err := json.Marshal(inputData)
	if err != nil {
		return nil, err
	}

	// TODO: the following ListenResult should be destroy after result is received
	tags := []string{uuid.NewV4().String()}
	resultC, errC, err := c.ListenResult(id, taskKey, "", tags)
	if err != nil {
		return nil, err
	}

	if _, err := c.client.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
		ServiceID:     id,
		TaskKey:       taskKey,
		InputData:     string(data),
		ExecutionTags: tags,
	}); err != nil {
		return nil, err
	}

	select {
	case r := <-resultC:
		if r.Error != "" {
			return nil, errors.New(r.Error)
		}
		return r, nil
	case err := <-errC:
		return nil, err
	}
}
