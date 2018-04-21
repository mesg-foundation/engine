package event

import (
	"fmt"

	"github.com/mesg-foundation/application/types"
)

// Listen
func (s *Server) Listen(request *types.ListenEventRequest, stream types.Event_ListenServer) (err error) {
	// service := service.New(request.Service)
	// stream.Send()
	fmt.Println("receive listen request")
	return
}
