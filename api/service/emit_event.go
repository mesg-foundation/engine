package service

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/event"
)

// EmitEvent permits to send and event to anyone who subscribed to it
func (s *Server) EmitEvent(context context.Context, request *EmitEventRequest) (reply *EmitEventReply, err error) {
	service, err := services.Get(request.ServiceHash)
	if err != nil {
		return
	}
	var data interface{}
	err = json.Unmarshal([]byte(request.EventData), &data)
	if err != nil {
		return
	}
	event, err := event.Create(&service, request.EventKey, data)
	if err != nil {
		return
	}
	event.Publish()
	reply = &EmitEventReply{}
	return
}
