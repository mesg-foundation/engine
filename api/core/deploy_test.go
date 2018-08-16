package core

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/cnf/structhash"
	"github.com/logrusorgru/aurora"
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
	require.True(t, stringSliceContains(stream.statuses, fmt.Sprintf("%s Completed.", aurora.Green("✔"))))
}

// TODO(ilgooz) also add tests for receiving chunks.
type testDeployStream struct {
	url       string // Git repo url.
	err       error
	serviceID string
	statuses  []string
	grpc.ServerStream
}

func newTestDeployStream(url string) *testDeployStream {
	return &testDeployStream{url: url}
}

func (s *testDeployStream) Send(m *DeployServiceReply) error {
	s.serviceID = m.GetServiceID()
	s.statuses = append(s.statuses, m.GetStatus())
	return nil
}

func (s *testDeployStream) Recv() (*DeployServiceRequest, error) {
	return &DeployServiceRequest{
		Value: &DeployServiceRequest_Url{Url: s.url},
	}, s.err
}

func stringSliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
