package service

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/event"
)

func (s *Server) EmitEvent(context context.Context, request *EmitEventRequest) (reply *EmitEventReply, err error) {
	service := request.Service
	var data interface{}
	err = json.Unmarshal([]byte(request.EventData), &data)
	if err != nil {
		return
	}
	event, err := event.Create(service, request.EventKey, data)
	if err != nil {
		return
	}
	event.Publish()
	reply = &EmitEventReply{}
	return
}
