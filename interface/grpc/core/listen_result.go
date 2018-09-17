package core

import (
	"encoding/json"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/protobuf/core"
)

// ListenResult listens for results from a services.
func (s *Server) ListenResult(request *core.ListenResultRequest, stream core.Core_ListenResultServer) error {
	ln, err := s.api.ListenResult(request.ServiceID,
		api.ListenResultTaskFilter(request.TaskFilter),
		api.ListenResultOutputFilter(request.OutputFilter),
		api.ListenResultTagFilters(request.TagFilters))
	if err != nil {
		return err
	}
	defer ln.Close()

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

			if err := stream.Send(&core.ResultData{
				ExecutionID:   execution.ID,
				TaskKey:       execution.Task,
				OutputKey:     execution.Output,
				OutputData:    string(outputs),
				ExecutionTags: execution.Tags,
			}); err != nil {
				return err
			}
		}
	}
}
