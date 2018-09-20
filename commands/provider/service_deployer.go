package provider

import (
	"context"
	"errors"
	"io"

	"github.com/asaskevich/govalidator"
	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// StatusType indicates the type of status message.
type StatusType int

const (
	// RUNNING indicates that status message belongs to an active state.
	RUNNING StatusType = iota + 1

	// DONE indicates that status message belongs to completed state.
	DONE
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
func (p *ServiceProvider) ServiceDeploy(path string, statuses chan DeployStatus) (id string,
	validationError, err error) {
	stream, err := p.client.DeployService(context.Background())
	if err != nil {
		return "", nil, err
	}

	deployment := make(chan deploymentResult)
	go readDeployReply(stream, deployment, statuses)

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

	if err := stream.CloseSend(); err != nil {
		return "", nil, err
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

func readDeployReply(stream coreapi.Core_DeployServiceClient, deployment chan deploymentResult,
	statuses chan DeployStatus) {
	result := deploymentResult{}

	for {
		message, err := stream.Recv()
		if err != nil {
			result.err = err
			deployment <- result
			return
		}

		var (
			status          = message.GetStatus()
			serviceID       = message.GetServiceID()
			validationError = message.GetValidationError()
		)

		switch {
		case status != nil:
			switch status.Type {
			case coreapi.DeployServiceReply_Status_RUNNING:
				statuses <- DeployStatus{
					Message: status.Message,
					Type:    RUNNING,
				}

			case coreapi.DeployServiceReply_Status_DONE:
				statuses <- DeployStatus{
					Message: status.Message,
					Type:    DONE,
				}
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
