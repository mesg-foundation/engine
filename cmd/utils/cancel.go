package utils

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitForCancel create a chan than is resolved when the user press CTRL+C
func WaitForCancel() chan os.Signal {
	abort := make(chan os.Signal, 1)
	signal.Notify(abort, syscall.SIGINT, syscall.SIGTERM)
	return abort
}
