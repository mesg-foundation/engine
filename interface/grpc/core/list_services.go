package core

import (
	"context"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	service "github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/x/xerrors"
)

// ListServices lists services.
func (s *Server) ListServices(ctx context.Context, request *coreapi.ListServicesRequest) (*coreapi.ListServicesReply, error) {
	var filters []api.ListServicesFilter
	if request.FilterActive {
		filters = append(filters, api.ListRunningServicesFilter())
	}
	services, err := s.api.ListServices(filters...)
	if err != nil {
		return nil, err
	}

	protoServices := toProtoServices(services)

	var (
		servicesLen = len(services)
		errC        = make(chan error, servicesLen)
	)

	// fill services status info.
	for _, s := range services {
		go func(s *service.Service) {
			status, err := s.Status()
			if err == nil {
				for _, ss := range protoServices {
					if ss.ID == s.ID {
						ss.Status = toProtoServiceStatusType(status)
					}
				}
			}
			errC <- err
		}(s)

	}

	var errs xerrors.Errors

	for i := 0; i < servicesLen; i++ {
		if err := <-errC; err != nil {
			errs = append(errs, err)
		}
	}

	return &coreapi.ListServicesReply{Services: protoServices}, errs.ErrorOrNil()
}
