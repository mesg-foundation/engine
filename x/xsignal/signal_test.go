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

package xsignal

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

// WaitForInterrupt creates a read channel for catch SIGINT and SIGTERM signals.
func TestWaitForInterrupt(t *testing.T) {
	abort := WaitForInterrupt()
	defer signal.Stop(abort)

	syscall.Kill(os.Getpid(), syscall.SIGINT)
	select {
	case sig := <-abort:
		if sig != syscall.SIGINT {
			t.Fatalf("unexpected signal - got %v, want %v", sig, syscall.SIGINT)
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for signal")
	}

	syscall.Kill(os.Getpid(), syscall.SIGTERM)
	select {
	case sig := <-abort:
		if sig != syscall.SIGTERM {
			t.Fatalf("unexpected signal - got %v, want %v", sig, syscall.SIGTERM)
		}
	case <-time.After(1 * time.Second):
		t.Fatalf("timeout waiting for signal")
	}
}
