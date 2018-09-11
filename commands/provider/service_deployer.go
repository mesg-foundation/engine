package provider

import (
	"context"
	"io"

	"github.com/asaskevich/govalidator"
	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/interface/grpc/core"
)

func (p *ServiceProvider) ServiceDeploy(path string) (string, bool, error) {
	stream, err := p.client.DeployService(context.Background())
	if err != nil {
		return "", false, err
	}

	deployment := make(chan deploymentResult)
	go readDeployReply(stream, deployment)

	if govalidator.IsURL(path) {
		if err := stream.Send(&core.DeployServiceRequest{
			Value: &core.DeployServiceRequest_Url{Url: path},
		}); err != nil {
			return "", true, err
		}
	} else {
		if err := deployServiceSendServiceContext(path, stream); err != nil {
			return "", true, err
		}
	}

	if err := stream.CloseSend(); err != nil {
		return "", true, err
	}

	result := <-deployment
	return result.serviceID, result.isValid, result.err
}

func deployServiceSendServiceContext(path string, stream core.Core_DeployServiceClient) error {
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

		if err := stream.Send(&core.DeployServiceRequest{
			Value: &core.DeployServiceRequest_Chunk{Chunk: buf[:n]},
		}); err != nil {
			return err
		}
	}

	return nil
}

type deploymentResult struct {
	serviceID string
	err       error
	isValid   bool
}

func readDeployReply(stream core.Core_DeployServiceClient, deployment chan deploymentResult) {
	var (
		result = deploymentResult{isValid: true}
	)

	for {
		message, err := stream.Recv()
		if err != nil {
			result.err = err
			deployment <- result
			return
		}

		var (
			serviceID       = message.GetServiceID()
			validationError = message.GetValidationError()
		)

		switch {
		case serviceID != "":
			result.serviceID = serviceID
			deployment <- result
			return

		case validationError != "":
			result.isValid = false
			deployment <- result
			return
		}
	}
}
