package xnet

import (
	"net"
	"testing"
)

func TestJoinHostPort(t *testing.T) {
	for _, tt := range []struct {
		host     string
		port     int
		hostPort string
	}{
		{"localhost", 80, "localhost:80"},
	} {
		if hostPort := JoinHostPort(tt.host, tt.port); hostPort != tt.hostPort {
			t.Errorf("JoinHostPort(%q, %q) = %q; want %q", tt.host, tt.port, hostPort, tt.hostPort)
		}
	}
}

func TestSplitHostPort(t *testing.T) {
	for _, tt := range []struct {
		hostPort string
		host     string
		port     int
	}{
		{"localhost:80", "localhost", 80},
	} {
		if host, port, err := SplitHostPort(tt.hostPort); host != tt.host || port != tt.port || err != nil {
			t.Errorf("SplitHostPort(%q) = %q, %q, %v; want %q, %q, nil", tt.hostPort, host, port, err, tt.host, tt.port)
		}
	}

	for _, tt := range []struct {
		hostPort string
		err      string
	}{
		{"localhost", "missing port in address"},
		{"localhost:a", "can't parse port"},
		{"localhost:-1", "port out of range"},
		{"localhost:65536", "port out of range"},
	} {
		if host, port, err := SplitHostPort(tt.hostPort); err == nil {
			t.Errorf("SplitHostPort(%q) should have failed", tt.hostPort)
		} else {
			e := err.(*net.AddrError)
			if e.Err != tt.err {
				t.Errorf("SplitHostPort(%q) = _, _, %q; want %q", tt.hostPort, e.Err, tt.err)
			}
			if host != "" || port != 0 {
				t.Errorf("SplitHostPort(%q) = %q, %q, err; want %q, %q, err on failure", tt.hostPort, host, port, "", "")
			}
		}
	}
}
