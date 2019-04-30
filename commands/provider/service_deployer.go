package provider

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/docker/cli/cli/command/image/build"
	"github.com/docker/docker/pkg/archive"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	validator "gopkg.in/go-playground/validator.v9"
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
	sid             string
	hash            string
	err             error
	validationError error
}

// ServiceDeploy deploys service from given path.
func (p *ServiceProvider) ServiceDeploy(path string, env map[string]string, statuses chan DeployStatus) (sid string, hash string, validationError, err error) {
	stream, err := p.client.DeployService(context.Background())
	if err != nil {
		return "", "", nil, err
	}

	deployment := make(chan deploymentResult)
	go func() {
		<-stream.Context().Done()
		deployment <- deploymentResult{err: stream.Context().Err()}
	}()
	go readDeployReply(stream, deployment, statuses)

	isURL := validator.New().Var(path, "url")
	switch {
	case strings.HasPrefix(path, "mesg:"):
		err = p.deployServiceFromMarketplace(path, env, stream)
	case isURL == nil:
		err = stream.Send(&coreapi.DeployServiceRequest{
			Value: &coreapi.DeployServiceRequest_Url{Url: path},
			Env:   env,
		})
	default:
		err = deployServiceSendServiceContext(path, env, stream)
	}
	if err != nil {
		return "", "", nil, err
	}

	if err := stream.CloseSend(); err != nil {
		return "", "", nil, err
	}

	result := <-deployment
	close(statuses)
	return result.sid, result.hash, result.validationError, result.err
}

func (p *ServiceProvider) deployServiceFromMarketplace(u string, env map[string]string, stream coreapi.Core_DeployServiceClient) error {

	urlParsed, err := url.Parse(u)
	if err != nil {
		return err
	}
	path := strings.Split(strings.Trim(urlParsed.EscapedPath(), "/"), "/")
	if urlParsed.Hostname() != "marketplace" || len(path) != 2 || path[0] != "service" || len(path[1]) == 0 {
		return fmt.Errorf("marketplace url %s invalid", u)
	}

	// Get ALL address from wallet
	addresses, err := p.wp.List()
	if err != nil {
		return err
	}

	// Check if one of them are is authorized
	authorized, sid, source, sourceType, err := p.mp.IsAuthorized("", path[1], addresses)
	if err != nil {
		return err
	}

	if !authorized {
		return fmt.Errorf("you are not authorized to deploy this service. Did you purchase it?\nExecute the following command to purchase it:\n\tmesg-core marketplace purchase %s", sid)
	}

	var url string
	switch sourceType {
	case "https", "http":
		url = source
	case "ipfs":
		url = "https://gateway.ipfs.io/ipfs/" + source
	default:
		return fmt.Errorf("unknown protocol %s", sourceType)
	}

	return stream.Send(&coreapi.DeployServiceRequest{
		Value: &coreapi.DeployServiceRequest_Url{Url: url},
		Env:   env,
	})
}

func deployServiceSendServiceContext(path string, env map[string]string, stream coreapi.Core_DeployServiceClient) error {
	excludes, err := build.ReadDockerignore(path)
	if err != nil {
		return err
	}

	if err := build.ValidateContextDirectory(path, excludes); err != nil {
		return fmt.Errorf("error checking context: %s", err)
	}

	// And canonicalize dockerfile name to a platform-independent one
	excludes = build.TrimBuildFilesFromExcludes(excludes, build.DefaultDockerfileName, false)

	archive, err := archive.TarWithOptions(path, &archive.TarOptions{
		Compression:     archive.Gzip,
		ExcludePatterns: excludes,
	})
	if err != nil {
		return err
	}

	if len(env) > 0 {
		if err := stream.Send(&coreapi.DeployServiceRequest{Env: env}); err != nil {
			return err
		}
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
			service         = message.GetService()
			validationError = message.GetValidationError()
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

		case service != nil:
			result.sid = service.Sid
			result.hash = service.Hash
			deployment <- result
			return

		case validationError != "":
			result.validationError = errors.New(validationError)
			deployment <- result
			return
		}
	}
}
