package event

import (
	"context"
	"log"

	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/pubsub"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/types"
)

// Emit a new event
func (s *Server) Emit(context context.Context, request *types.EmitEventRequest) (reply *types.EventReply, err error) {
	log.Println("Receive event", aurora.Green(request.Event), "from", aurora.Green(request.Service.Name), ":", aurora.Bold(request.Data))

	service := service.New(request.Service)

	reply = &types.EventReply{
		Event: request.Event,
		Data:  request.Data,
	}

	go pubsub.Publish(service.EventSubscriptionChannel(), reply)

	return
}
