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

package acknowledgement

import (
	fmt "fmt"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	statusKey   = "status"
	statusReady = "ready"
)

// WaitForStreamToBeReady waits until the server publish a header with status ready on the stream.
func WaitForStreamToBeReady(stream grpc.ClientStream) error {
	// wait for header.
	header, err := stream.Header()
	if err != nil {
		return err
	}
	// header received. check status.
	statuses := header.Get(statusKey)
	statusesLen := len(statuses)
	if statusesLen == 0 {
		return fmt.Errorf("stream header does not contain any status")
	}
	lastStatus := statuses[statusesLen-1]
	if lastStatus != statusReady {
		return fmt.Errorf("stream header status is different than ready. Got %q", lastStatus)
	}
	return nil
}

// SetStreamReady send an header on the stream to notify the server is ready.
func SetStreamReady(stream grpc.ServerStream) error {
	return stream.SendHeader(metadata.Pairs(statusKey, statusReady))
}
