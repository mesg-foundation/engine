package core

import (
	"github.com/mesg-foundation/core/mesg"
	service "github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/importer"
)

// DeployService saves a service in the database and returns the hash of this service.
func (s *Server) DeployService(stream Core_DeployServiceServer) error {
	statuses := make(chan string, 0)
	go sendDeployStatus(statuses, stream)

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
	if url != "" {
		service, validationError, err = s.mesg.DeployServiceFromURL(url, mesg.DeployServiceStatusOption(statuses))
	} else {
		service, validationError, err = s.mesg.DeployService(sr, mesg.DeployServiceStatusOption(statuses))
	}

	if err != nil {
		return err
	}
	if validationError != nil {
		return stream.Send(&DeployServiceReply{
			Value: &DeployServiceReply_ValidationError{ValidationError: validationError.Error()},
		})
	}

	return stream.Send(&DeployServiceReply{
		Value: &DeployServiceReply_ServiceID{ServiceID: service.Id},
	})
}

func sendDeployStatus(statuses chan string, stream Core_DeployServiceServer) {
	for status := range statuses {
		stream.Send(&DeployServiceReply{
			Value: &DeployServiceReply_Status{Status: status},
		})
	}
}

type deployServiceStreamReader struct {
	stream Core_DeployServiceServer

	data []byte
	i    int64
}

func newDeployServiceStreamReader(stream Core_DeployServiceServer) *deployServiceStreamReader {
	return &deployServiceStreamReader{stream: stream}
}

func (r *deployServiceStreamReader) GetURL() (url string, err error) {
	message, err := r.stream.Recv()
	r.data = message.GetChunk()
	return message.GetUrl(), err
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
