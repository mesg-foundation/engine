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
	syscall.Kill(os.Getpid(), syscall.SIGTERM)

	sigCount := 0

loop:
	for {
		select {
		case sig := <-abort:
			if sig != syscall.SIGTERM && sig != syscall.SIGINT {
				t.Fatalf("unexpected signal - got %v", sig)
			}
			if sigCount = sigCount + 1; sigCount >= 2 {
				break loop
			}
		case <-time.After(1 * time.Second):
			t.Fatalf("timeout waiting for signal")
		}
	}

	if sigCount != 2 {
		t.Fatalf("received signals failed: got: %d, want: %d", sigCount, 2)
	}
}
