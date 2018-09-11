package clierrors

import (
	"fmt"

	"github.com/docker/docker/client"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	cannotReachTheCore = "Cannot reach the Core"
	startCore          = "Please start the core by running: mesg-core start"
	cannotReachDocker  = "Cannot reach Docker"
	installDocker      = "Please make sure Docker is running.\nIf Docker is not installed on your machine you can install it here: https://store.docker.com/search?type=edition&offering=community"
)

func ErrorMessage(err error) string {
	switch {
	case coreConnectionError(err):
		return fmt.Sprintf("%s\n%s", cannotReachTheCore, startCore)
	case dockerDaemonError(err):
		return fmt.Sprintf("%s\n%s", cannotReachDocker, installDocker)
	default:
		return err.Error()
	}
}

func coreConnectionError(err error) bool {
	s := status.Convert(err)
	return s.Code() == codes.Unavailable
}

func dockerDaemonError(err error) bool {
	return client.IsErrConnectionFailed(err)
}
