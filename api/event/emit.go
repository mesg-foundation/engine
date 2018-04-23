package event

import (
	"context"

	"github.com/mesg-foundation/application/pubsub"
	"github.com/mesg-foundation/application/service"
	"github.com/mesg-foundation/application/types"
)

// Emit a new event
func (s *Server) Emit(context context.Context, request *types.EmitEventRequest) (reply *types.EventReply, err error) {

	service := service.New(request.Service)

	reply = &types.EventReply{
		Data:  request.Data,
		Event: request.Event,
	}

	go pubsub.Publish(service.EventSubscriptionKey(), reply)

	return
}
