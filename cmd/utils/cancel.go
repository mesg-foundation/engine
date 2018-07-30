package utils

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitForCancel creates a chan that is resolved when the user press CTRL+C.
func WaitForCancel() chan os.Signal {
	abort := make(chan os.Signal, 1)
	signal.Notify(abort, syscall.SIGINT, syscall.SIGTERM)
	return abort
}
