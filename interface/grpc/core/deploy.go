package core

import (
	"errors"
	"sync"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	service "github.com/mesg-foundation/core/service"
	"github.com/mesg-foundation/core/service/importer"
)

// DeployService deploys a service from Git URL or service.tar.gz file. It'll send status
// events during the process and finish with sending service id or validation error.
func (s *Server) DeployService(stream coreapi.Core_DeployServiceServer) error {
	var (
		statuses = make(chan api.DeployStatus)
		option   = api.DeployServiceStatusOption(statuses)
		wg       sync.WaitGroup

		service         *service.Service
		validationError *importer.ValidationError
		err             error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		sendDeployStatus(statuses, stream)
	}()

	// read first requesest from stream and check if it's url or tarball
	in, err := stream.Recv()
	if err != nil {
		return err
	}

	// env must be set with first package (always)
	env := in.GetEnv()

	if url := in.GetUrl(); url != "" {
		service, validationError, err = s.api.DeployServiceFromURL(url, env, option)
	} else {
		// create tarball reader with first chunk of bytes
		tarball := &deployChunkReader{
			stream: stream,
			buf:    in.GetChunk(),
		}
		service, validationError, err = s.api.DeployService(tarball, env, option)
	}

	// wait for statuses to be sent first, otherwise sending multiple messages at the
	// same time may cause messages to be sent in different order.
	wg.Wait()

	if err != nil {
		return err
	}
	if validationError != nil {
		return stream.Send(&coreapi.DeployServiceReply{
			Value: &coreapi.DeployServiceReply_ValidationError{
				ValidationError: validationError.Error(),
			},
		})
	}

	return stream.Send(&coreapi.DeployServiceReply{
		Value: &coreapi.DeployServiceReply_Service_{
			Service: &coreapi.DeployServiceReply_Service{
				Sid:  service.Sid,
				Hash: service.Hash,
			},
		},
	})
}

func sendDeployStatus(statuses chan api.DeployStatus, stream coreapi.Core_DeployServiceServer) {
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

// deployChunkReader implements io.Reader for stream chunks.
type deployChunkReader struct {
	stream coreapi.Core_DeployServiceServer

	buf []byte
	i   int
}

func (r *deployChunkReader) Read(p []byte) (n int, err error) {
	if r.i >= len(r.buf) {
		in, err := r.stream.Recv()
		if err != nil {
			return 0, err
		}

		r.buf, r.i = in.GetChunk(), 0
		if len(r.buf) == 0 {
			return 0, errors.New("deploy: got empty chunk of tarball")
		}
	}
	n = copy(p, r.buf[r.i:])
	r.i += n
	return n, nil
}
