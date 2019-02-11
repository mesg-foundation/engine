package servicetest

import (
	"net"

	"github.com/mesg-foundation/core/protobuf/serviceapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

// Socket is an in-memory gRPC socket.
type Socket struct {
	ln *bufconn.Listener
}

func newSocket() *Socket {
	s := &Socket{
		ln: bufconn.Listen(1024),
	}
	return s
}

// Dial returns the client side net conn.
func (s *Socket) Dial() (net.Conn, error) {
	return s.ln.Dial()
}

func (s *Socket) listen(serviceServer serviceapi.ServiceServer) error {
	server := grpc.NewServer()
	serviceapi.RegisterServiceServer(server, serviceServer)
	return server.Serve(s.ln)
}

func (s *Socket) close() error {
	return s.ln.Close()
}
