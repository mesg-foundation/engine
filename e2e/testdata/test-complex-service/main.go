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
	"github.com/mesg-foundation/engine/server/grpc/runner"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	// env variables for configure mesg client.
	envMesgEndpoint  = "MESG_ENDPOINT"
	envMesgMsg       = "MESG_MSG"
	envMesgSignature = "MESG_SIGNATURE"
)

// newClient creates a new client from env variables supplied by mesg engine.
func newClient() (runner.RunnerClient, error) {
	endpoint := os.Getenv(envMesgEndpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("env %q is empty", envMesgEndpoint)
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

	return runner.NewRunnerClient(conn), nil
}

func register(client runner.RunnerClient) (string, error) {
	msg := os.Getenv(envMesgMsg)
	if msg == "" {
		return "", fmt.Errorf("env %q is empty", envMesgMsg)
	}
	signature := os.Getenv(envMesgSignature)
	if signature == "" {
		return "", fmt.Errorf("env %q is empty", envMesgSignature)
	}

	resp, err := client.Register(context.Background(), &runner.RegisterRequest{
		Msg:       msg,
		Signature: signature,
	})
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

// sendEvent creates a new event.
func sendEvent(client runner.RunnerClient, token string, key string) {
	log.Println("sending event:", key)
	if _, err := client.Event(context.Background(), &runner.EventRequest{
		Key: key,
	}, grpc.PerRPCCredentials(runner.NewTokenCredential(token))); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	client, err := newClient()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected to %s\n", os.Getenv(envMesgEndpoint))

	token, err := register(client)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered with token %s\n", token)

	// give some time to nginx to start
	time.Sleep(10 * time.Second)

	sendEvent(client, token, "test_service_ready")

	// check env default value
	if os.Getenv("ENVA") == "do_not_override" {
		sendEvent(client, token, "read_env_ok")
	} else {
		sendEvent(client, token, "read_env_error")
	}

	// check env override value
	if os.Getenv("ENVB") == "is_override" {
		sendEvent(client, token, "read_env_override_ok")
	} else {
		sendEvent(client, token, "read_env_override_error")
	}

	if err := ioutil.WriteFile("/volume/test/test.txt", []byte("foo"), 0644); err == nil {
		sendEvent(client, token, "access_volumes_ok")
	} else {
		sendEvent(client, token, "access_volumes_error")
	}

	if _, err := http.Get("http://nginx:80/"); err == nil {
		sendEvent(client, token, "resolve_dependence_ok")
	} else {
		sendEvent(client, token, "resolve_dependence_error")
	}

	if content, err := ioutil.ReadFile("/etc/nginx/nginx.conf"); len(content) > 0 && err == nil {
		sendEvent(client, token, "access_volumes_from_ok")
	} else {
		sendEvent(client, token, "access_volumes_from_error")
	}

	<-xsignal.WaitForInterrupt()
}
