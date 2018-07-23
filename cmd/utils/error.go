package utils

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return aurora.Sprintf("%s\nPlease start the core by running: mesg-core start", aurora.Red("Cannot reach the Core."))
	default:
		return aurora.Red(err.Error()).String()
	}
}

func coreConnectionError(err error) bool {
	s := status.Convert(err)
	return s.Code() == codes.Unavailable
}
