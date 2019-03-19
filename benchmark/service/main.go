package main

import (
	"context"
	"log"
	"os"

	"github.com/mesg-foundation/core/protobuf/serviceapi"
	"google.golang.org/grpc"
)

func main() {
	endpoint := os.Getenv("MESG_ENDPOINT")
	token := os.Getenv("MESG_TOKEN")

	sh := &statsHandler{}

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure(), grpc.WithStatsHandler(sh))
	if err != nil {
		log.Fatal(err)
	}

	client := serviceapi.NewServiceClient(conn)

	l, err := client.ListenTask(context.Background(), &serviceapi.ListenTaskRequest{Token: token})
	if err != nil {
		log.Fatal(err)
	}

	sh.print()

	for {
		if _, err := l.Recv(); err != nil {
			log.Fatal(err)
		}
	}
}
