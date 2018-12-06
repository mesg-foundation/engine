package core

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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
		statuses      = make(chan api.DeployStatus)
		urls          = make(chan string)
		confirmations = make(chan bool)
		r, w          = io.Pipe()
		wg            sync.WaitGroup
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		sendDeployStatus(statuses, stream)
	}()

	var (
		service         *service.Service
		validationError *importer.ValidationError
		err             error
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := readStream(stream, urls, confirmations, w)
		fmt.Println("panic error deploy")
		panic(err)
	}()
	// sr := newDeployServiceStreamReader(stream, confirmations)
	// url, err := sr.GetURL()
	// if err != nil {
	// return err
	// }
	deployOptions := []api.DeployServiceOption{
		api.DeployServiceStatusOption(statuses),
		api.DeployServiceConfirmationsOption(confirmations),
	}
	// if url != "" {
	if url := <-urls; url != "" {
		service, validationError, err = s.api.DeployServiceFromURL(url, deployOptions...)
	} else {
		service, validationError, err = s.api.DeployService(r, deployOptions...)
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
		case api.Confirmation:
			typ = coreapi.DeployServiceReply_Status_CONFIRMATION
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

func readStream(stream coreapi.Core_DeployServiceServer, urls chan string, confirmations chan bool, w io.WriteCloser) error {
	for {
		fmt.Println("read stream for")
		message, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		fmt.Println("check message type")

		value := message.GetValue()
		if x, ok := value.(*coreapi.DeployServiceRequest_Chunk); ok {
			fmt.Println("is chunk")
			if bytes.Equal(x.Chunk, []byte("END_OF_SERVICE")) { //TODO: improve END_OF_SERVICE code. Try to use a standard code.
				if err := w.Close(); err != nil {
					return err
				}
			}
			if _, err := w.Write(x.Chunk); err != nil {
				return errors.New("error on write. " + err.Error())
			}
			fmt.Println("chunk written")
			urls <- ""
		} else if x, ok := value.(*coreapi.DeployServiceRequest_Confirmation); ok {
			fmt.Println("is confirmation")
			confirmations <- x.Confirmation
		} else if x, ok := value.(*coreapi.DeployServiceRequest_Url); ok {
			fmt.Println("is url")
			urls <- x.Url
		} else {
			return errors.New("unknown type")
		}

	}
}

// type deployServiceStreamReader struct {
// 	stream        coreapi.Core_DeployServiceServer
// 	confirmations chan bool

// 	data []byte
// 	i    int64
// }

// func newDeployServiceStreamReader(stream coreapi.Core_DeployServiceServer, confirmations chan bool) *deployServiceStreamReader {
// 	return &deployServiceStreamReader{
// 		stream:        stream,
// 		confirmations: confirmations,
// 	}
// }

// func (r *deployServiceStreamReader) GetURL() (url string, err error) {
// 	message, err := r.stream.Recv()
// 	if err != nil {
// 		return "", err
// 	}
// 	r.data = message.GetChunk()
// 	return message.GetUrl(), err
// }

// func (r *deployServiceStreamReader) GetConfirmation() (bool, error) {
// 	message, err := r.stream.Recv()
// 	if err != nil {
// 		return false, err
// 	}
// 	r.data = message.GetChunk()
// 	return message.GetConfirmation(), err
// }

// func (r *deployServiceStreamReader) Read(p []byte) (n int, err error) {
// 	if r.i >= int64(len(r.data)) {
// 		message, err := r.stream.Recv()
// 		if err != nil {
// 			return 0, err
// 		}

// 		value := message.GetValue()
// 		// if x, ok := value.(*coreapi.DeployServiceRequest_Url); ok {
// 		// 	// x.Url
// 		// } else
// 		if x, ok := value.(*coreapi.DeployServiceRequest_Chunk); ok {
// 			if bytes.Equal(x.Chunk, []byte("END_OF_SERVICE")) {
// 				return 0, io.EOF
// 			}
// 			r.data = x.Chunk
// 			r.i = 0
// 		} else if x, ok := value.(*coreapi.DeployServiceRequest_Confirmation); ok {
// 			r.confirmations <- x.Confirmation
// 		}
// 		return r.Read(p)
// 	}
// 	n = copy(p, r.data[r.i:])
// 	r.i += int64(n)
// 	return n, nil
// }
