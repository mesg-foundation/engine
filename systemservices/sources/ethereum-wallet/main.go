package main

import (
	"context"
	"encoding/json"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/mesg-foundation/core/protobuf/serviceapi"
	"google.golang.org/grpc"
)

type eventData struct {
	Content string `json:"content"`
}

func main() {

	client, err := initServiceClient()
	if err != nil {
		panic(err)
	}

	ks := keystore.NewKeyStore("/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	account, err := ks.NewAccount("pass")
	if err != nil {
		panic(err)
	}
	err = emitEvent(client, account.Address.String())
	if err != nil {
		panic(err)
	}

	// service, err := mesg.New()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println("Service is started")
	// err = service.Listen(
	// 	mesg.Task("resolve", resolveHandler),
	// 	mesg.Task("addPeers", addPeersHandler),
	// )
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func initServiceClient() (serviceapi.ServiceClient, error) {
	conn, err := grpc.DialContext(context.Background(), os.Getenv("MESG_ENDPOINT"), grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return serviceapi.NewServiceClient(conn), nil
}

func emitEvent(client serviceapi.ServiceClient, content string) error {
	dataBytes, err := json.Marshal(&eventData{
		Content: content,
	})
	if err != nil {
		return err
	}
	_, err = client.EmitEvent(context.Background(), &serviceapi.EmitEventRequest{
		Token:     os.Getenv("MESG_TOKEN"),
		EventKey:  "test",
		EventData: string(dataBytes),
	})
	if err != nil {
		return err
	}
	return nil
}
