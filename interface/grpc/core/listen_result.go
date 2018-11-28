package core

import (
	"encoding/json"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// ListenResult listens for results from a services.
func (s *Server) ListenResult(request *coreapi.ListenResultRequest, stream coreapi.Core_ListenResultServer) error {
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

			data := &coreapi.ResultData{
				ExecutionID:   execution.ID,
				TaskKey:       execution.TaskKey,
				OutputKey:     execution.OutputKey,
				OutputData:    string(outputs),
				ExecutionTags: execution.Tags,
			}
			if execution.Error != nil {
				data.Result = &coreapi.ResultData_Error_{
					Error: &coreapi.ResultData_Error{
						Message: execution.Error.Error(),
					},
				}
			} else {
				data.Result = &coreapi.ResultData_Output_{
					Output: &coreapi.ResultData_Output{
						OutputKey:  execution.OutputKey,
						OutputData: string(outputs),
					},
				}
			}
			if err := stream.Send(data); err != nil {
				return err
			}
		}
	}
}
