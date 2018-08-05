// TODO(ilgooz): remove this file after service package made Newable.
package daemon

import (
	"log"

	"github.com/mesg-foundation/core/container"
)

var defaultContainer *container.Container

func init() {
	c, err := container.New()
	if err != nil {
		log.Fatal(err)
	}
	defaultContainer = c
}
