package service

import (
	"fmt"

	types "github.com/mesg-foundation/application/types"
	"golang.org/x/net/context"
)

func (s *Server) Stop(ctx context.Context, request *types.StopServiceRequest) (reply *types.ServiceReply, err error) {
	fmt.Println("stop")
	return
}
