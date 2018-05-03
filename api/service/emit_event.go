package service

import (
	"context"

	"github.com/mesg-foundation/core/pubsub"
)

// EmitEvent a new event
func (s *Server) EmitEvent(context context.Context, request *EmitEventRequest) (reply *EmitEventReply, err error) {
	channel := request.Service.EventSubscriptionChannel()

	reply = &EmitEventReply{}

	go pubsub.Publish(channel, reply)

	return
}
