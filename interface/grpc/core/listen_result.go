package core

import (
	"encoding/json"
	"time"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"google.golang.org/grpc/metadata"
)

// ListenResult listens for results from a services.
func (s *Server) ListenResult(request *coreapi.ListenResultRequest, stream coreapi.Core_ListenResultServer) error {
	// FOR TEST ONLY: pause the server to simulate bad condition
	time.Sleep(3 * time.Second)

	ln, err := s.api.ListenResult(request.ServiceID,
		api.ListenResultTaskFilter(request.TaskFilter),
		api.ListenResultOutputFilter(request.OutputFilter),
		api.ListenResultTagFilters(request.TagFilters))
	if err != nil {
		return err
	}
	defer ln.Close()
	// send header to confirm to client that the server is ready
	err = stream.SendHeader(metadata.Pairs("status", "ready"))
	if err != nil {
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
			outputs, err := json.Marshal(execution.OutputData)
			if err != nil {
				return err
			}
			if err := stream.Send(&coreapi.ResultData{
				ExecutionID:   execution.ID,
				TaskKey:       execution.TaskKey,
				OutputKey:     execution.OutputKey,
				OutputData:    string(outputs),
				ExecutionTags: execution.Tags,
				Error:         execution.Error,
			}); err != nil {
				return err
			}
		}
	}
}
