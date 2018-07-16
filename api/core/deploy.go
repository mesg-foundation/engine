package core

import (
	"context"
	"github.com/mesg-foundation/core/database/services"
)

// DeployService save a service in the database and return the hash of this service
func (s *Server) DeployService(ctx context.Context, request *DeployServiceRequest) (reply *DeployServiceReply, err error) {
	service := request.Service
	hash, err := services.Save(service)
	if err != nil {
		return
	}
	reply = &DeployServiceReply{
		ServiceID: hash,
	}
	return
}
