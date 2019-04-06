package core

import (
	"encoding/json"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/protobuf/acknowledgement"
	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// ListenResult listens for results from a services.
func (s *Server) ListenResult(request *coreapi.ListenResultRequest, stream coreapi.Core_ListenResultServer) error {
	ln, err := s.api.ListenResult(request.ServiceID,
		api.ListenResultTaskFilter(request.TaskFilter),
		api.ListenResultTagFilters(request.TagFilters))
	if err != nil {
		return err
	}
	defer ln.Close()

	// send header to notify client that the stream is ready.
	if err := acknowledgement.SetStreamReady(stream); err != nil {
		return err
	}

	ctx := stream.Context()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case err := <-ln.Err:
			return err

		case execution := <-ln.Executions:
			outputs, err := json.Marshal(execution.Outputs)
			if err != nil {
				return err
			}
			if err := stream.Send(&coreapi.ResultData{
				ExecutionID:   execution.ID,
				TaskKey:       execution.TaskKey,
				Data:          string(outputs),
				ExecutionTags: execution.Tags,
				Error:         execution.Error,
			}); err != nil {
				return err
			}
		}
	}
}
