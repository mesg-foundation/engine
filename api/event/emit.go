package event

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mesg-foundation/application/types"
)

// Emit
func (s *Server) Emit(context context.Context, request *types.EmitEventRequest) (reply *types.EventReply, err error) {
	// service := service.New(request.Service)
	// stream.Send()
	log.Println("receive emit request")
	log.Println("Event", request.Event)
	log.Println("Data", request.Data)

	// data := &Data{}
	var data interface{}
	json.Unmarshal([]byte(request.Data), &data)
	decoded := data.(map[string]interface{})
	log.Println("data.number", decoded["number"].(float64))

	reply = &types.EventReply{
		Event: "ok",
	}

	return
}
