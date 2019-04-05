package core

import (
	"context"
	"sync"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/x/xerrors"
)

// ListServices lists services.
func (s *Server) ListServices(ctx context.Context, request *coreapi.ListServicesRequest) (*coreapi.ListServicesReply, error) {
	services, err := s.api.ListServices()
	if err != nil {
		return nil, err
	}

	var (
		protoServices []*coreapi.Service
		mp            sync.Mutex

		servicesLen = len(services)
		errC        = make(chan error, servicesLen)
		wg          sync.WaitGroup
	)

	wg.Add(servicesLen)
	for _, ss := range services {
		go func(ss *service.Service) {
			defer wg.Done()
			status, err := ss.Status(s.container)
			if err != nil {
				errC <- err
				return
			}
			protoService := toProtoService(ss)
			protoService.Status = toProtoServiceStatusType(status)
			mp.Lock()
			protoServices = append(protoServices, protoService)
			mp.Unlock()
		}(ss)
	}

	wg.Wait()
	close(errC)

	var errs xerrors.Errors
	for err := range errC {
		errs = append(errs, err)
	}

	return &coreapi.ListServicesReply{Services: protoServices}, errs.ErrorOrNil()
}
