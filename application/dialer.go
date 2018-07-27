package application

import (
	"net"
	"time"
)

// Dialer describes a dialer implementation for network.
type Dialer interface {
	Dial() (net.Conn, error)
}

type gRPCDialer struct {
	dialer Dialer
}

func newGRPCDialer(dialer Dialer) *gRPCDialer {
	return &gRPCDialer{dialer}
}

// Dial used to produce a net.Conn to connect gRPC server with.
func (d *gRPCDialer) Dial(addr string, timeout time.Duration) (net.Conn, error) {
	return d.dialer.Dial()
}
