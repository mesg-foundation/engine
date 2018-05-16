package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/mesg-foundation/core/api/core"
	"github.com/mesg-foundation/core/service"
	"google.golang.org/grpc"
)

var cli core.CoreClient

func when(service *service.Service, event string, then func(data *core.EventData)) {
	stream, err := cli.ListenEvent(context.Background(), &core.ListenEventRequest{
		Service: service,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("waiting for event to trigger")
	for {
		e, _ := stream.Recv()
		if strings.Compare(e.EventKey, event) == 0 {
			then(e)
		}
	}
}

func then(service *service.Service, task string, inputs func(data *core.EventData) interface{}) func(data *core.EventData) {
	return func(eventData *core.EventData) {
		in := inputs(eventData)
		d, err := json.Marshal(in)
		if err != nil {
			log.Fatalln(err)
		}
		_, err = cli.ExecuteTask(context.Background(), &core.ExecuteTaskRequest{
			Service:  service,
			TaskData: string(d),
			TaskKey:  task,
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func start(path string) (*service.Service, func()) {
	service, err := service.ImportFromPath(path)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = cli.StartService(context.Background(), &core.StartServiceRequest{
		Service: service,
	})
	if err != nil {
		log.Fatalln(err)
	}
	return service, func() {
		cli.StopService(context.Background(), &core.StopServiceRequest{Service: service})
	}
}

func main() {
	connection, _ := grpc.Dial(":50052", grpc.WithInsecure())
	cli = core.NewCoreClient(connection)

	devcon, stopDevcon := start("../../services/devcon-update")
	defer stopDevcon()
	slack, stopSlack := start("../../services/slack")
	defer stopSlack()

	when(devcon, "update",
		then(slack, "notify", func(data *core.EventData) interface{} {
			return map[string]string{
				"channel": "general",
				"title":   "Update on https://devcon.ethereum.org",
				"token":   os.Getenv("SLACK_TOKEN"),
			}
		}))
}
