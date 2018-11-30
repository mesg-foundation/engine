package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mesg-foundation/core/protobuf/serviceapi"
)

// SubmitResult submits results of an execution.
func (s *Server) SubmitResult(context context.Context, request *serviceapi.SubmitResultRequest) (*serviceapi.SubmitResultReply, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(request.OutputData), &data); err != nil {
		return nil, fmt.Errorf("service sent invalid json data in output %q", request.OutputKey)
	}
	return &serviceapi.SubmitResultReply{}, s.api.SubmitResult(request.ExecutionID, request.OutputKey, data)
}
