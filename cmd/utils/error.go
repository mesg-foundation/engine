package utils

import (
	"fmt"
	"os"

	"github.com/docker/docker/client"
	"github.com/logrusorgru/aurora"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	cannotReachTheCore = "Cannot reach the Core."
	startCore          = "Please start the core by running: mesg-core start"
	cannotReachDocker  = "Please make sure Docker is running"
	installDocker      = "If Docker is not installed on your machine you can install it here https://store.docker.com/search?type=edition&offering=community"
)

// HandleError display the error and stop the process if error exist
func HandleError(err error) {
	if err != nil {
		fmt.Println(errorMessage(err))
		os.Exit(0)
	}
}

func errorMessage(err error) string {
	switch {
	case coreConnectionError(err):
		return aurora.Sprintf("%s\n%s", aurora.Red(cannotReachTheCore), startCore)
	case dockerDaemonError(err):
		return aurora.Sprintf("%s\n%s", aurora.Red(cannotReachDocker), installDocker)
	default:
		return aurora.Red(err.Error()).String()
	}
}

func coreConnectionError(err error) bool {
	s := status.Convert(err)
	return s.Code() == codes.Unavailable
}

func dockerDaemonError(err error) bool {
	return client.IsErrConnectionFailed(err)
}
