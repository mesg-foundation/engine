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

package core

import (
	"context"
	"testing"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/stretchr/testify/require"
)

func TestListServices(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"
	server, closer := newServer(t)
	defer closer()

	stream := newTestDeployStream(url)
	require.NoError(t, server.DeployService(stream))
	defer server.api.DeleteService(stream.serviceID, false)

	reply, err := server.ListServices(context.Background(), &coreapi.ListServicesRequest{})
	require.NoError(t, err)

	services, err := server.api.ListServices()
	require.NoError(t, err)

	apiProtoServices := toProtoServices(services)

	require.Len(t, apiProtoServices, 1)
	require.Equal(t, reply.Services[0].Hash, apiProtoServices[0].Hash)
}
