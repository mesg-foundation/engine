package result

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/mesg-foundation/core/execution"
	"github.com/mesg-foundation/core/types"
)

// Submit a task result
func (s *Server) Submit(context context.Context, request *types.SubmitResultRequest) (reply *types.ResultReply, err error) {
	execution := execution.InProgress(request.ExecutionID)
	if execution == nil {
		err = errors.New("No task in progress with the ID " + request.ExecutionID)
		return
	}
	var data interface{}
	err = json.Unmarshal([]byte(request.Data), &data)
	if err != nil {
		return
	}
	reply, err = execution.Complete(request.Output, data)
	return
}
