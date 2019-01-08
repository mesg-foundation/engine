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

// +build integration

package core

import (
	"testing"

	"github.com/mesg-foundation/core/api"
	"github.com/stretchr/testify/require"
)

func TestIntegrationDeployService(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"

	server, closer := newServer(t)
	defer closer()
	stream := newTestDeployStream(url)

	require.Nil(t, server.DeployService(stream))
	defer server.api.DeleteService(stream.serviceID, false)

	require.Len(t, stream.serviceID, 40)
	require.Contains(t, stream.statuses, api.DeployStatus{
		Message: "Image built with success",
		Type:    api.DonePositive,
	})
}
