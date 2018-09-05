package core

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/cnf/structhash"
	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/api"
	"github.com/stretchr/testify/require"
	grpc "google.golang.org/grpc"
)

func TestDeployService(t *testing.T) {
	url := "https://github.com/mesg-foundation/service-webhook"

	server, dt := newServerAndDockerTest(t)
	dt.ProvideImageBuild(ioutil.NopCloser(strings.NewReader(`{"stream":"sha256:x"}`)), nil)

	stream := newTestDeployStream(url)
	require.Nil(t, server.DeployService(stream))
	require.Equal(t, 1, structhash.Version(stream.serviceID))

	require.Contains(t, stream.statuses, api.DeployStatus{
		Message: fmt.Sprintf("%s Completed.", aurora.Green("✔")),
		Type:    api.DONE,
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

func (s *testDeployStream) Send(m *DeployServiceReply) error {
	s.serviceID = m.GetServiceID()

	status := m.GetStatus()
	if status != nil {
		var typ api.StatusType
		switch status.Type {
		case DeployServiceReply_Status_RUNNING:
			typ = api.RUNNING
		case DeployServiceReply_Status_DONE:
			typ = api.DONE
		}
		s.statuses = append(s.statuses, api.DeployStatus{
			Message: status.Message,
			Type:    typ,
		})
	}

	return nil
}

func (s *testDeployStream) Recv() (*DeployServiceRequest, error) {
	return &DeployServiceRequest{
		Value: &DeployServiceRequest_Url{Url: s.url},
	}, s.err
}
