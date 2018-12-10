package provider

import (
	"context"
	"errors"
	"io"

	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/asaskevich/govalidator"
	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// StatusType indicates the type of status message.
type StatusType int

const (
	_ StatusType = iota // skip zero value.

	// Running indicates that status message belongs to a continuous state.
	Running

	// DonePositive indicates that status message belongs to a positive noncontinuous state.
	DonePositive

	// DoneNegative indicates that status message belongs to a negative noncontinuous state.
	DoneNegative
)

// DeployStatus represents the deployment status.
type DeployStatus struct {
	Message string
	Type    StatusType
}

// deploymentResult keeps information about deployment result.
type deploymentResult struct {
	serviceID       string
	err             error
	validationError error
}

// ServiceDeploy deploys service from given path.
func (p *ServiceProvider) ServiceDeploy(path string, confirmation *bool,
	confirmationFunc func(sid string) (deletion bool, err error),
	statuses chan DeployStatus) (id string,
	validationError, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := p.client.DeployService(ctx)
	if err != nil {
		return "", nil, err
	}
	defer stream.CloseSend()

	deployment := make(chan deploymentResult)
	go readDeployReply(stream, deployment, statuses, confirmationFunc)

	if confirmation != nil {
		if err := stream.Send(&coreapi.DeployServiceRequest{
			Value: &coreapi.DeployServiceRequest_Confirmation{
				Confirmation: &wrappers.BoolValue{Value: *confirmation},
			},
		}); err != nil {
			return "", nil, err
		}
	}

	if govalidator.IsURL(path) {
		if err := stream.Send(&coreapi.DeployServiceRequest{
			Value: &coreapi.DeployServiceRequest_Url{Url: path},
		}); err != nil {
			return "", nil, err
		}
	} else {
		if err := deployServiceSendServiceContext(path, stream); err != nil {
			return "", nil, err
		}
	}

	result := <-deployment
	close(statuses)
	return result.serviceID, result.validationError, result.err
}

func deployServiceSendServiceContext(path string, stream coreapi.Core_DeployServiceClient) error {
	archive, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression: archive.Gzip,
	})
	if err != nil {
		return err
	}

	buf := make([]byte, 1024)
	for {
		n, err := archive.Read(buf)
		if err == io.EOF {
			if err := stream.Send(&coreapi.DeployServiceRequest{
				Value: &coreapi.DeployServiceRequest_ChunkDone{ChunkDone: true},
			}); err != nil {
				return err
			}
			break
		}
		if err != nil {
			return err
		}

		if err := stream.Send(&coreapi.DeployServiceRequest{
			Value: &coreapi.DeployServiceRequest_Chunk{Chunk: buf[:n]},
		}); err != nil {
			return err
		}
	}

	return nil
}

func readDeployReply(stream coreapi.Core_DeployServiceClient,
	deployment chan deploymentResult,
	statuses chan DeployStatus,
	confirmationFunc func(string) (bool, error)) {
	result := deploymentResult{}

	for {
		message, err := stream.Recv()
		if err != nil {
			result.err = err
			deployment <- result
			return
		}

		var (
			status              = message.GetStatus()
			requestConfirmation = message.GetRequestConfirmation()
			serviceID           = message.GetServiceID()
			validationError     = message.GetValidationError()
		)

		switch {
		case status != nil:
			s := DeployStatus{
				Message: status.Message,
			}

			switch status.Type {
			case coreapi.DeployServiceReply_Status_RUNNING:
				s.Type = Running
			case coreapi.DeployServiceReply_Status_DONE_POSITIVE:
				s.Type = DonePositive
			case coreapi.DeployServiceReply_Status_DONE_NEGATIVE:
				s.Type = DoneNegative
			}

			statuses <- s

		case requestConfirmation != "":
			var deletion bool
			if confirmationFunc != nil {
				var err error
				deletion, err = confirmationFunc(requestConfirmation)
				if err != nil {
					result.err = err
					deployment <- result
					return
				}
			}

			if err := stream.Send(&coreapi.DeployServiceRequest{
				Value: &coreapi.DeployServiceRequest_Confirmation{
					Confirmation: &wrappers.BoolValue{Value: deletion},
				},
			}); err != nil {
				result.err = err
				deployment <- result
				return
			}

		case serviceID != "":
			result.serviceID = serviceID
			deployment <- result
			return

		case validationError != "":
			result.validationError = errors.New(validationError)
			deployment <- result
			return
		}
	}
}
