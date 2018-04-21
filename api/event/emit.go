package event

import (
	"context"
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
	return
}
