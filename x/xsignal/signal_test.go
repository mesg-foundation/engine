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
