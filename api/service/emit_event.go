package service

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/database/services"
	"github.com/mesg-foundation/core/event"
)

// EmitEvent permits to send and event to anyone who subscribed to it
func (s *Server) EmitEvent(context context.Context, request *EmitEventRequest) (*EmitEventReply, error) {
	service, err := services.Get(request.Token)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal([]byte(request.EventData), &data)
	if err != nil {
		return nil, err
	}
	event, err := event.Create(&service, request.EventKey, data)
	if err != nil {
		return nil, err
	}
	event.Publish()
	return &EmitEventReply{}, nil
}
