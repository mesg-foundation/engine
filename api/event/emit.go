package event

import (
	"context"

	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
)

// Emit a new event
func (s *Server) Emit(context context.Context, request *types.EmitEventRequest) (reply *types.EventReply, err error) {
	channel := service.New(request.Service).EventSubscriptionChannel()

	reply = &types.EventReply{
		Event: request.Event,
		Data:  request.Data,
	}

	go pubsub.Publish(channel, reply)

	return
}
