// Copyright 2018 MESG Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	installDocker      = "Please make sure Docker is running\nIf Docker is not installed on your machine you can install it here: https://store.docker.com/search?type=edition&offering=community"
)

// ErrorMessage returns error description based on error type.
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
