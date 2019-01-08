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

	"github.com/mesg-foundation/core/protobuf/coreapi"
)

// GetService returns service serviceID.
func (s *Server) GetService(ctx context.Context, request *coreapi.GetServiceRequest) (*coreapi.GetServiceReply, error) {
	ss, err := s.api.GetService(request.ServiceID)
	if err != nil {
		return nil, err
	}
	protoService := toProtoService(ss)
	status, err := ss.Status()
	if err != nil {
		return nil, err
	}
	protoService.Status = toProtoServiceStatusType(status)
	return &coreapi.GetServiceReply{Service: protoService}, nil
}
