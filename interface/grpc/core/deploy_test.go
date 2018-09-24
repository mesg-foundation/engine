package core

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/mesg-foundation/core/api"
	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/stretchr/testify/require"
	grpc "google.golang.org/grpc"
)

func TestDeployService(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"

	server, dt := newServerAndDockerTest(t)
	dt.ProvideImageBuild(ioutil.NopCloser(strings.NewReader(`{"stream":"sha256:x"}`)), nil)

	stream := newTestDeployStream(url)
	require.Nil(t, server.DeployService(stream))
	require.Len(t, stream.serviceID, 40)

	require.Contains(t, stream.statuses, api.DeployStatus{
		Message: "Service deployed.",
		Type:    api.DonePositive,
	})
}

// TODO(ilgooz) also add tests for receiving chunks.
type testDeployStream struct {
	url       string // Git repo url.
	err       error
	serviceID string
	statuses  []api.DeployStatus
	grpc.ServerStream
}

func newTestDeployStream(url string) *testDeployStream {
	return &testDeployStream{url: url}
}

func (s *testDeployStream) Send(m *coreapi.DeployServiceReply) error {
	s.serviceID = m.GetServiceID()

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
