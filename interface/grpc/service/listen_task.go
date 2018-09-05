package service

import (
	"encoding/json"
)

// ListenTask creates a stream that will send data for every task to execute.
func (s *Server) ListenTask(request *ListenTaskRequest, stream Service_ListenTaskServer) error {
	ln, err := s.api.ListenTask(request.Token)
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
			inputs, err := json.Marshal(execution.Inputs)
			if err != nil {
				return err
			}

			if err := stream.Send(&TaskData{
				ExecutionID: execution.ID,
				TaskKey:     execution.Task,
				InputData:   string(inputs),
			}); err != nil {
				return err
			}
		}
	}
}
