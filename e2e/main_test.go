package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/mesg-foundation/engine/core"
	"github.com/mesg-foundation/engine/cosmos"
	"github.com/mesg-foundation/engine/hash"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var testkb *cosmos.Keybase

var passmd = metadata.Pairs(
	"credential_username", "engine",
	"credential_passphrase", "password",
)

type apiclient struct {
	pb.ServiceClient
	pb.EventClient
	pb.ExecutionClient
	pb.AccountClient
	pb.ProcessClient
	pb.InstanceClient
	pb.OwnershipClient
}

var client apiclient

func newapiclient() apiclient {
	conn, err := grpc.DialContext(context.Background(), "localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return apiclient{
		pb.NewServiceClient(conn),
		pb.NewEventClient(conn),
		pb.NewExecutionClient(conn),
		pb.NewAccountClient(conn),
		pb.NewProcessClient(conn),
		pb.NewInstanceClient(conn),
		pb.NewOwnershipClient(conn),
	}
}

func readCreateServiceRequest(filename string) *pb.CreateServiceRequest {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var req pb.CreateServiceRequest
	if err = json.Unmarshal(b, &req); err != nil {
		log.Fatal(err)
	}
	return &req
}

var (
	testServiceHash  hash.Hash
	testInstanceHash hash.Hash
)

func TestMain(m *testing.M) {
	//os.Setenv("MESG_LOG_LEVEL", "fatal")

	kb, cleanup := core.Start()
	testkb = kb

	client = newapiclient()

	// deploy services for test
	resp, err := client.ServiceClient.Create(
		metadata.NewOutgoingContext(context.Background(), passmd),
		readCreateServiceRequest("testdata/test-service/compile.json"),
	)
	if err != nil {
		log.Fatal(err)
	}
	testServiceHash = resp.Hash

	iresp, err := client.InstanceClient.Create(
		metadata.NewOutgoingContext(context.Background(), passmd),
		&pb.CreateInstanceRequest{ServiceHash: testServiceHash},
	)
	if err != nil {
		log.Fatal(err)
	}
	testInstanceHash = iresp.Hash

	time.Sleep(1 * time.Second)

	ret := m.Run()
	cleanup()
	os.Exit(ret)
}
