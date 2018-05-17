package core

import (
	"errors"

	"github.com/mesg-foundation/core/cmd/daemon"
	"golang.org/x/net/context"
)

func getDeamonIP() (daemonIP string, err error) {
	daemonContainer, err := daemon.Container()
	if err != nil {
		return
	}
	if daemonContainer == nil {
		err = errors.New("Daemon container is not found")
		return
	}
	networkContainer := daemonContainer.Networks.Networks["mesg-shared-network"]
	if networkContainer.IPAddress == "" {
		err = errors.New("Network 'mesg-shared-network' not found")
		return
	}
	daemonIP = networkContainer.IPAddress
	return
}

var sharedNetwork = "mesg-shared-network"

// Start a service
func (s *Server) StartService(ctx context.Context, request *StartServiceRequest) (reply *StartServiceReply, err error) {
	service := request.Service
	daemonIP, err := getDeamonIP()
	if err != nil {
		return
	}	
	_, err = service.Start(daemonIP, sharedNetwork)
	reply = &StartServiceReply{}
	return
}
