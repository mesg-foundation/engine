package core

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"github.com/mesg-foundation/core/sdk"
	"github.com/stretchr/testify/require"
	grpc "google.golang.org/grpc"
)

func TestDeployService(t *testing.T) {
	url := "git://github.com/mesg-foundation/service-webhook"

	server, dt, closer := newServerAndDockerTest(t)
	defer closer()
	dt.ProvideImageBuild(ioutil.NopCloser(strings.NewReader(`{"stream":"sha256:x"}`)), nil)

	stream := newTestDeployStream(url)
	require.Nil(t, server.DeployService(stream))
	require.Len(t, stream.sid, 7)
	require.NotEmpty(t, stream.hash)

	require.Contains(t, stream.statuses, sdk.DeployStatus{
		Message: "Image built with success",
		Type:    sdk.DonePositive,
	})
}

// TODO(ilgooz) also add tests for receiving chunks.
type testDeployStream struct {
	url      string // Git repo url.
	err      error
	sid      string
	hash     string
	statuses []sdk.DeployStatus
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
		var typ sdk.StatusType
		switch status.Type {
		case coreapi.DeployServiceReply_Status_RUNNING:
			typ = sdk.Running
		case coreapi.DeployServiceReply_Status_DONE_POSITIVE:
			typ = sdk.DonePositive
		case coreapi.DeployServiceReply_Status_DONE_NEGATIVE:
			typ = sdk.DoneNegative
		}
		s.statuses = append(s.statuses, sdk.DeployStatus{
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
