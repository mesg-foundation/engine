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

package main

import (
	"context"
	"time"

	"github.com/mesg-foundation/core/systemservices/sources/resolver/proto/core"
	"github.com/mesg-foundation/core/x/xstrings"
	mesg "github.com/mesg-foundation/go-service"
	"google.golang.org/grpc"
)

const (
	resolveFoundOutputsKey    = "found"
	resolveNotFoundOutputsKey = "notFound"
)

type resolveInputs struct {
	ServiceID string `json:"serviceID"`
}

type resolveFoundOutputs struct {
	Address   string `json:"address"`
	ServiceID string `json:"serviceID"`
}

type resolveNotFoundOutputs struct {
	ServiceID string `json:"serviceID"`
}

func resolveHandler(execution *mesg.Execution) (string, mesg.Data) {
	var inputs resolveInputs
	if err := execution.Data(&inputs); err != nil {
		return newOutputsError(err)
	}
	for _, peer := range getPeers() {
		match, err := matchingPeer(peer, inputs.ServiceID)
		if err != nil {
			return newOutputsError(err)
		}
		if match {
			return resolveFoundOutputsKey, &resolveFoundOutputs{
				Address:   peer,
				ServiceID: inputs.ServiceID,
			}
		}
	}
	return resolveNotFoundOutputsKey, &resolveNotFoundOutputs{
		ServiceID: inputs.ServiceID,
	}
}

func matchingPeer(peer, serviceID string) (bool, error) {
	serviceIDs, err := listServices(peer)
	if err != nil {
		return false, err
	}
	return xstrings.SliceContains(serviceIDs, serviceID), nil
}

func listServices(peer string) ([]string, error) {
	conn, client, err := newCoreClient(peer)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	reply, err := client.ListServices(ctx, &core.ListServicesRequest{})
	if err != nil {
		return nil, err
	}
	serviceIDs := make([]string, len(reply.Services))
	for i, service := range reply.Services {
		serviceIDs[i] = service.ID
	}
	return serviceIDs, nil
}

func newCoreClient(address string) (*grpc.ClientConn, core.CoreClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure())
	if err != nil {
		return nil, nil, err
	}
	client := core.NewCoreClient(conn)
	return conn, client, nil
}
