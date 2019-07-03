package client

import (
	"net"
	"os"
	"testing"

	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/server/grpc/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestClient(t *testing.T) {
	l, err := net.Listen("tcp", ":50052")
	require.NoError(t, err)

	s := grpc.NewServer()
	pb.RegisterEventServer(s, api.NewEventServer(nil))
	pb.RegisterExecutionServer(s, api.NewExecutionServer(nil))
	pb.RegisterInstanceServer(s, api.NewInstanceServer(nil))
	pb.RegisterServiceServer(s, api.NewServiceServer(nil))
	defer s.Stop()

	go func() {
		require.NoError(t, s.Serve(l))
	}()

	_, err = New()
	require.Contains(t, err.Error(), envMesgEndpoint)
	os.Setenv(envMesgEndpoint, l.Addr().String())

	_, err = New()
	require.Contains(t, err.Error(), envMesgInstanceHash)
	os.Setenv(envMesgInstanceHash, "1")

	_, err = New()
	require.NoError(t, err)
}
