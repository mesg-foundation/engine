package service

import (
	"net"
	"time"
)

type Dialer interface {
	Dial() (net.Conn, error)
}

type gRPCDialer struct {
	dialer Dialer
}

func newGRPCDialer(dialer Dialer) *gRPCDialer {
	return &gRPCDialer{dialer}
}

func (d *gRPCDialer) Dial(addr string, timeout time.Duration) (net.Conn, error) {
	return d.dialer.Dial()
}
