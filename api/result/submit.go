package result

import (
	"context"

	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
)

// Submit a task result
func (s *Server) Submit(context context.Context, request *types.SubmitResultRequest) (reply *types.ResultReply, err error) {
	service := service.New(request.Service)

	reply = &types.ResultReply{
		Task:   request.Task,
		Output: request.Output,
		Data:   request.Data,
	}

	go pubsub.Publish(service.ResultSubscriptionChannel(), reply)

	return
}
