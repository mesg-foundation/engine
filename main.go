package main

import (
	"github.com/mesg-foundation/engine/core"
	"github.com/mesg-foundation/engine/x/xsignal"
)

func main() {
	_, cleanup := core.Start()
	defer cleanup()
	<-xsignal.WaitForInterrupt()
}
