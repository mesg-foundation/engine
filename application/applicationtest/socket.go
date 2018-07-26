package applicationtest

import (
	"net"

	"github.com/mesg-foundation/core/api/core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type Socket struct {
	ln *bufconn.Listener
}

func newSocket() *Socket {
	s := &Socket{
		ln: bufconn.Listen(1024),
	}
	return s
}

func (s *Socket) Dial() (net.Conn, error) {
	return s.ln.Dial()
}

func (s *Socket) listen(coreServer *coreServer) error {
	server := grpc.NewServer()
	core.RegisterCoreServer(server, coreServer)
	return server.Serve(s.ln)
}

func (s *Socket) close() error {
	return s.ln.Close()
}
