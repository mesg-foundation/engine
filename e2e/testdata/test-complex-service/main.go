package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/mesg-foundation/engine/ext/xsignal"
	"github.com/mesg-foundation/engine/hash"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	// env variables for configure mesg client.
	envMesgEndpoint     = "MESG_ENDPOINT"
	envMesgInstanceHash = "MESG_INSTANCE_HASH"
)

// Client is a client to connect to all mesg exposed API.
type Client struct {
	// all clients registered by mesg server.
	pb.EventClient
	pb.ExecutionClient

	// instance hash that could be used in api calls.
	InstanceHash hash.Hash
}

// New creates a new client from env variables supplied by mesg engine.
func New() (*Client, error) {
	endpoint := os.Getenv(envMesgEndpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("client: mesg server address env(%s) is empty", envMesgEndpoint)
	}

	instanceHash, err := hash.Decode(os.Getenv(envMesgInstanceHash))
	if err != nil {
		return nil, fmt.Errorf("client: error with mesg's instance hash env(%s): %s", envMesgInstanceHash, err.Error())
	}

	dialoptions := []grpc.DialOption{
		// Keep alive prevents Docker network to drop TCP idle connections after 15 minutes.
		// See: https://forum.mesg.com/t/solution-summary-for-docker-dropping-connections-after-15-min/246
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: 5 * time.Minute, // 5 minutes is the minimun time of gRPC enforcement policy.
		}),
		grpc.WithTimeout(10 * time.Second),
		grpc.WithInsecure(),
	}

	conn, err := grpc.DialContext(context.Background(), endpoint, dialoptions...)
	if err != nil {
		return nil, fmt.Errorf("client: connection error: %s", err)
	}

	return &Client{
		ExecutionClient: pb.NewExecutionClient(conn),
		EventClient:     pb.NewEventClient(conn),
		InstanceHash:    instanceHash,
	}, nil
}

// SendEvent creates a new event.
func (c *Client) SendEvent(key string) {
	if _, err := c.EventClient.Create(context.Background(), &pb.CreateEventRequest{
		InstanceHash: c.InstanceHash,
		Key:          key,
	}); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client, err := New()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected to %s\n", os.Getenv(envMesgEndpoint))

	// give some time to nginx to start
	time.Sleep(10 * time.Second)

	client.SendEvent("test_service_ready")

	// check env default value
	if os.Getenv("ENVA") == "do_not_override" {
		client.SendEvent("read_env_ok")
	} else {
		client.SendEvent("read_env_error")
	}

	// check env override value
	if os.Getenv("ENVB") == "is_override" {
		client.SendEvent("read_env_override_ok")
	} else {
		client.SendEvent("read_env_override_error")
	}

	if err := ioutil.WriteFile("/volume/test/test.txt", []byte("foo"), 0644); err == nil {
		client.SendEvent("access_volumes_ok")
	} else {
		client.SendEvent("access_volumes_error")
	}

	if _, err := http.Get("http://nginx:80/"); err == nil {
		client.SendEvent("resolve_dependence_ok")
	} else {
		client.SendEvent("resolve_dependence_error")
	}

	if content, err := ioutil.ReadFile("/etc/nginx/nginx.conf"); len(content) > 0 && err == nil {
		client.SendEvent("access_volumes_from_ok")
	} else {
		client.SendEvent("access_volumes_from_error")
	}

	<-xsignal.WaitForInterrupt()
}
