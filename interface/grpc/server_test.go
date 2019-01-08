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

package grpc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	waitForServe = 500 * time.Millisecond
)

func TestServerServe(t *testing.T) {
	s := New("localhost:50052", nil, nil)
	go func() {
		time.Sleep(waitForServe)
		s.Close()
	}()
	err := s.Serve()
	require.NoError(t, err)
}

func TestServerServeNoAddress(t *testing.T) {
	s := Server{}
	go func() {
		time.Sleep(waitForServe)
		s.Close()
	}()
	err := s.Serve()
	require.Error(t, err)
}

func TestServerListenAfterClose(t *testing.T) {
	s := New("localhost:50052", nil, nil)
	go s.Serve()
	time.Sleep(waitForServe)
	s.Close()
	require.Equal(t, &alreadyClosedError{}, s.Serve())
}
