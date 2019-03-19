package main

import (
	"context"
	"log"
	"os"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"google.golang.org/grpc"
)

func main() {
	token := os.Getenv("MESG_TOKEN")
	if token == "" {
		token = "benchmark-service"
	}

	endpoint := os.Getenv("MESG_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:50052"
	}

	sh := &statsHandler{}

	conn, err := grpc.Dial(endpoint, grpc.WithInsecure(), grpc.WithStatsHandler(sh))
	if err != nil {
		log.Fatal(err)
	}

	client := coreapi.NewCoreClient(conn)

	l, err := client.ListenEvent(context.Background(), &coreapi.ListenEventRequest{
		ServiceID:   token,
		EventFilter: "foo",
	})
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
