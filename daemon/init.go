package daemon

import (
	"log"

	"github.com/mesg-foundation/core/container"
)

var defaultContainer container.Container

// TODO(ilgooz): remove init after daemon package made Newable.
func init() {
	c, err := container.New()
	if err != nil {
		log.Fatal(err)
	}
	defaultContainer = c
}
