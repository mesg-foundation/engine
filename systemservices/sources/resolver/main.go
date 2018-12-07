package main

import (
	"fmt"
	"log"

	mesg "github.com/mesg-foundation/go-service"
)

func main() {
	service, err := mesg.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Service is started")
	err = service.Listen(
		mesg.Task("resolve", resolveHandler),
		mesg.Task("addPeers", addPeersHandler),
	)
	if err != nil {
		log.Fatal(err)
	}
}
