package service

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/protobuf/service"
)

// EmitEvent permits to send and event to anyone who subscribed to it.
func (s *Server) EmitEvent(context context.Context, request *service.EmitEventRequest) (*service.EmitEventReply, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(request.EventData), &data); err != nil {
		return nil, err
	}
	return &service.EmitEventReply{}, s.api.EmitEvent(request.Token, request.EventKey, data)
}
