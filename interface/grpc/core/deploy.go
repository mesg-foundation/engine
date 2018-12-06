package core

import (
	"sync"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	service "github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/importer"
)

// DeployService deploys a service from Git URL or service.tar.gz file. It'll send status
// events during the process and finish with sending service id or validation error.
// TODO(ilgooz): sync `stream.Send()`s by doing it in a single goroutine.
func (s *Server) DeployService(stream coreapi.Core_DeployServiceServer) error {
	var (
		statuses = make(chan api.DeployStatus)
		wg       sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		forwardDeployStatuses(statuses, stream)
	}()

	var (
		service         *service.Service
		validationError *importer.ValidationError
		err             error
	)

	sr := newDeployServiceStreamReader(stream)
	url, err := sr.GetURL()
	if err != nil {
		return err
	}

	deployOptions := []api.DeployServiceOption{
		api.DeployServiceStatusOption(statuses),
		api.DeployServiceConfirmationOption(func(alias string) bool {
			// request for confirmation.
			if err := stream.Send(&coreapi.DeployServiceReply{
				Value: &coreapi.DeployServiceReply_RequestConfirmation{RequestConfirmation: alias},
			}); err != nil {
				return false
			}

			// receive the confirmation result.
			// TODO(ilgooz) add timeout.
			deletion, err := sr.GetConfirmation()
			if err != nil {
				return false
			}
			return deletion
		}),
	}

	if url != "" {
		service, validationError, err = s.api.DeployServiceFromURL(url, deployOptions...)
	} else {
		service, validationError, err = s.api.DeployService(sr, deployOptions...)
	}
	wg.Wait()

	if err != nil {
		return err
	}
	if validationError != nil {
		return stream.Send(&coreapi.DeployServiceReply{
			Value: &coreapi.DeployServiceReply_ValidationError{ValidationError: validationError.Error()},
		})
	}

	return stream.Send(&coreapi.DeployServiceReply{
		Value: &coreapi.DeployServiceReply_ServiceID{ServiceID: service.ID},
	})
}

func forwardDeployStatuses(statuses chan api.DeployStatus, stream coreapi.Core_DeployServiceServer) {
	for status := range statuses {
		var typ coreapi.DeployServiceReply_Status_Type
		switch status.Type {
		case api.Running:
			typ = coreapi.DeployServiceReply_Status_RUNNING
		case api.DonePositive:
			typ = coreapi.DeployServiceReply_Status_DONE_POSITIVE
		case api.DoneNegative:
			typ = coreapi.DeployServiceReply_Status_DONE_NEGATIVE
		}
		stream.Send(&coreapi.DeployServiceReply{
			Value: &coreapi.DeployServiceReply_Status_{
				Status: &coreapi.DeployServiceReply_Status{
					Message: status.Message,
					Type:    typ,
				},
			},
		})
	}
}

type deployServiceStreamReader struct {
	stream coreapi.Core_DeployServiceServer

	data []byte
	i    int64
}

func newDeployServiceStreamReader(stream coreapi.Core_DeployServiceServer) *deployServiceStreamReader {
	return &deployServiceStreamReader{
		stream: stream,
	}
}

func (r *deployServiceStreamReader) GetURL() (url string, err error) {
	message, err := r.stream.Recv()
	if err != nil {
		return "", err
	}
	r.data = message.GetChunk()
	return message.GetUrl(), err
}

func (r *deployServiceStreamReader) GetConfirmation() (bool, error) {
	message, err := r.stream.Recv()
	if err != nil {
		return false, err
	}
	return message.GetConfirmation(), err
}

func (r *deployServiceStreamReader) Read(p []byte) (n int, err error) {
	if r.i >= int64(len(r.data)) {
		message, err := r.stream.Recv()
		if err != nil {
			return 0, err
		}
		r.data = message.GetChunk()
		r.i = 0
		return r.Read(p)
	}
	n = copy(p, r.data[r.i:])
	r.i += int64(n)
	return n, nil
}
