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
	"github.com/mesg-foundation/engine/server/grpc/orchestrator"
	"github.com/mesg-foundation/engine/server/grpc/runner"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	// env variables for configure mesg client.
	envMesgEndpoint          = "MESG_ENDPOINT"
	envMesgServiceHash       = "MESG_SERVICE_HASH"
	envMesgEnvHash           = "MESG_ENV_HASH"
	envMesgRegisterSignature = "MESG_REGISTER_SIGNATURE"
)

// register
func register() (string, error) {
	endpoint := os.Getenv(envMesgEndpoint)
	if endpoint == "" {
		return "", fmt.Errorf("env %q is empty", envMesgEndpoint)
	}

	conn, err := grpc.DialContext(context.Background(), endpoint, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	client := orchestrator.NewRunnerClient(conn)

	serviceHash, err := hash.Decode(os.Getenv(envMesgServiceHash))
	if err != nil {
		return "", err
	}
	envHash, err := hash.Decode(os.Getenv(envMesgEnvHash))
	if err != nil {
		return "", err
	}
	signature := os.Getenv(envMesgRegisterSignature)

	resp, err := client.Register(context.Background(), &orchestrator.RunnerRegisterRequest{
		ServiceHash: serviceHash,
		EnvHash:     envHash,
	}, grpc.PerRPCCredentials(&signCred{signature}))
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func newClient(token string) (runner.RunnerClient, error) {
	endpoint := os.Getenv(envMesgEndpoint)
	if endpoint == "" {
		return nil, fmt.Errorf("env %q is empty", envMesgEndpoint)
	}

	// runner client
	dialoptions := []grpc.DialOption{
		// Keep alive prevents Docker network to drop TCP idle connections after 15 minutes.
		// See: https://forum.mesg.com/t/solution-summary-for-docker-dropping-connections-after-15-min/246
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time: 5 * time.Minute, // 5 minutes is the minimun time of gRPC enforcement policy.
		}),
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(&tokenCred{token}),
	}
	conn, err := grpc.DialContext(context.Background(), endpoint, dialoptions...)
	if err != nil {
		return nil, fmt.Errorf("connection error: %s", err)
	}

	return runner.NewRunnerClient(conn), nil
}

// sendEvent creates a new event.
func sendEvent(client runner.RunnerClient, token string, key string) {
	log.Println("sending event:", key)
	if _, err := client.Event(context.Background(), &runner.EventRequest{
		Key: key,
	}); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	token, err := register()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered with token %s\n", token)

	client, err := newClient(token)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected to %s\n", os.Getenv(envMesgEndpoint))

	// give some time to nginx to start
	time.Sleep(10 * time.Second)

	sendEvent(client, token, "service_ready")

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

// tokenCred is a structure that manage a token.
type tokenCred struct {
	token string
}

// GetRequestMetadata returns the metadata for the request.
func (c *tokenCred) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		runner.CredentialToken: c.token,
	}, nil
}

// RequireTransportSecurity tells if the transport should be secured.
func (c *tokenCred) RequireTransportSecurity() bool {
	return false
}

// signCred is a structure that manage a token.
type signCred struct {
	signature string
}

// GetRequestMetadata returns the metadata for the request.
func (c *signCred) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		orchestrator.RequestSignature: c.signature,
	}, nil
}

// RequireTransportSecurity tells if the transport should be secured.
func (c *signCred) RequireTransportSecurity() bool {
	return false
}
