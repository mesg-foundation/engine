package xsignal

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitForInterrupt creates a read channel for catch SIGINT and SIGTERM signals.
func WaitForInterrupt() chan os.Signal {
	abort := make(chan os.Signal, 1)
	signal.Notify(abort, os.Interrupt, syscall.SIGTERM)
	return abort
}
