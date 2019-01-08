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
	"strconv"
)

// Maximum and minimum ports number.
const (
	MinPort = 0
	MaxPort = 65535
)

// JoinHostPort combines host and port into a network address of the
// form "host:port" or, if host contains a colon or a percent sign,
// "[host]:port".
func JoinHostPort(host string, port int) string {
	return net.JoinHostPort(host, strconv.FormatInt(int64(port), 10))
}

// SplitHostPort splits a network address of the form "host:port",
// "host%zone:port", "[host]:port" or "[host%zone]:port" into host or
// host%zone and port.
func SplitHostPort(hostport string) (string, int, error) {
	host, portStr, err := net.SplitHostPort(hostport)
	if err != nil {
		return "", 0, err
	}

	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		return "", 0, &net.AddrError{Err: "can't parse port", Addr: hostport}
	}

	if port < MinPort || port > MaxPort {
		return "", 0, &net.AddrError{Err: "port out of range", Addr: hostport}
	}

	return host, int(port), nil
}
