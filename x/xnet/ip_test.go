// Copyright 2018 MESG Foundation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
