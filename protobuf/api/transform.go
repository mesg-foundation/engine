package api

import (
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
)

// TransformCreateReqToService returns service from create service request.
func TransformCreateReqToService(req *CreateServiceRequest) *service.Service {
	srv := &service.Service{
		Sid:           req.Sid,
		Name:          req.Name,
		Description:   req.Description,
		Configuration: req.Configuration,
		Tasks:         req.Tasks,
		Events:        req.Events,
		Dependencies:  req.Dependencies,
		Repository:    req.Repository,
		Source:        req.Source,
	}

	srv.Hash = hash.Dump(srv)
	return srv
}
