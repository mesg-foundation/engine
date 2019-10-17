package main

import (
	"flag"

	"github.com/mesg-foundation/engine/core"
	"github.com/mesg-foundation/engine/x/xsignal"
)

var configpath = flag.String("config", "", "set yaml config path")

func main() {
	flag.Parse()
	_, closer := core.Start(*configpath)
	defer closer()
	<-xsignal.WaitForInterrupt()
}
