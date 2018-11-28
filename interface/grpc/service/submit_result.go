package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mesg-foundation/core/protobuf/serviceapi"
	"github.com/mesg-foundation/core/x/xerrors"
)

// SubmitResult submits results of an execution.
func (s *Server) SubmitResult(context context.Context, request *serviceapi.SubmitResultRequest) (*serviceapi.SubmitResultReply, error) {
	var data map[string]interface{}
	var errs xerrors.Errors

	rerr := json.Unmarshal([]byte(request.OutputData), &data)
	if rerr != nil {
		errs = append(errs, fmt.Errorf("invalid output data error: %s", rerr))
	}
	if err := s.api.SubmitResult(request.ExecutionID, request.OutputKey, data, rerr); err != nil {
		errs = append(errs, err)
	}

	return &serviceapi.SubmitResultReply{}, errs.ErrorOrNil()
}
