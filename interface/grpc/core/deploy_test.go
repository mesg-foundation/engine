package core

import (
	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	grpc "google.golang.org/grpc"
)

// TODO(ilgooz) also add tests for receiving chunks.
type testDeployStream struct {
	url      string // Git repo url.
	err      error
	sid      string
	hash     string
	statuses []api.DeployStatus
	grpc.ServerStream
}

func newTestDeployStream(url string) *testDeployStream {
	return &testDeployStream{url: url}
}

func (s *testDeployStream) Send(m *coreapi.DeployServiceReply) error {
	s.sid = m.GetService().GetSid()
	s.hash = m.GetService().GetHash()

	status := m.GetStatus()
	if status != nil {
		var typ api.StatusType
		switch status.Type {
		case coreapi.DeployServiceReply_Status_RUNNING:
			typ = api.Running
		case coreapi.DeployServiceReply_Status_DONE_POSITIVE:
			typ = api.DonePositive
		case coreapi.DeployServiceReply_Status_DONE_NEGATIVE:
			typ = api.DoneNegative
		}
		s.statuses = append(s.statuses, api.DeployStatus{
			Message: status.Message,
			Type:    typ,
		})
	}

	return nil
}

func (s *testDeployStream) Recv() (*coreapi.DeployServiceRequest, error) {
	return &coreapi.DeployServiceRequest{
		Value: &coreapi.DeployServiceRequest_Url{Url: s.url},
	}, s.err
}
