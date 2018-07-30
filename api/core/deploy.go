package core

import (
	"context"

	"github.com/mesg-foundation/core/database/services"
)

// DeployService saves a service in the database and returns the hash of this service.
func (s *Server) DeployService(ctx context.Context, request *DeployServiceRequest) (*DeployServiceReply, error) {
	service := request.Service
	hash, err := services.Save(service)
	if err != nil {
		return nil, err
	}
	return &DeployServiceReply{
		ServiceID: hash,
	}, nil
}
