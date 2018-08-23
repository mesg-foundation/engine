package core

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/cnf/structhash"
	"github.com/logrusorgru/aurora"
	"github.com/mesg-foundation/core/mesg"
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

	require.Contains(t, stream.statuses, mesg.DeployStatus{
		Message: fmt.Sprintf("%s Completed.", aurora.Green("✔")),
		Type:    mesg.DONE,
	})
}

// TODO(ilgooz) also add tests for receiving chunks.
type testDeployStream struct {
	url       string // Git repo url.
	err       error
	serviceID string
	statuses  []mesg.DeployStatus
	grpc.ServerStream
}

func newTestDeployStream(url string) *testDeployStream {
	return &testDeployStream{url: url}
}

func (s *testDeployStream) Send(m *DeployServiceReply) error {
	s.serviceID = m.GetServiceID()

	status := m.GetStatus()
	if status != nil {
		var typ mesg.StatusType
		switch status.Type {
		case DeployServiceReply_Status_RUNNING:
			typ = mesg.RUNNING
		case DeployServiceReply_Status_DONE:
			typ = mesg.DONE
		}
		s.statuses = append(s.statuses, mesg.DeployStatus{
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
