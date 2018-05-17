package core

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/mesg-foundation/core/service"
	"google.golang.org/grpc"
)

type workflow struct {
	service   *service.Service
	event     string
	listeners []*task
}

type task struct {
	service *service.Service
	task    string
	convert func(data *EventData) interface{}
}

var once sync.Once
var workflows []*workflow
var cli CoreClient

func connect() {
	once.Do(func() {
		connection, _ := grpc.Dial(":50052", grpc.WithInsecure())
		cli = NewCoreClient(connection)
	})
}

func start(service *service.Service) {
	cli.StartService(context.Background(), &StartServiceRequest{
		Service: service,
	})
}

func StartWorkflow() {
	abort := make(chan os.Signal, 1)
	signal.Notify(abort, syscall.SIGINT, syscall.SIGTERM)
	<-abort
	for _, wf := range workflows {
		cli.StopService(context.Background(), &StopServiceRequest{Service: wf.service})
		for _, task := range wf.listeners {
			cli.StopService(context.Background(), &StopServiceRequest{Service: task.service})
		}
	}
}

func When(service *service.Service, event string) (wf *workflow) {
	connect()
	start(service)
	stream, err := cli.ListenEvent(context.Background(), &ListenEventRequest{
		Service: service,
	})
	if err != nil {
		log.Fatalln(err)
	}
	wf = &workflow{
		service: service,
	}
	go func() {
		for {
			e, _ := stream.Recv()
			if strings.Compare(e.EventKey, event) == 0 {
				wf.executeAll(e)
			}
		}
	}()

	wf = &workflow{
		service: service,
	}
	workflows = append(workflows, wf)
	return
}

func (workflow *workflow) Then(service *service.Service, taskName string, convert func(data *EventData) interface{}) {
	connect()
	start(service)
	workflow.listeners = append(workflow.listeners, &task{
		service: service,
		task:    taskName,
		convert: convert,
	})
}

func (workflow *workflow) executeAll(event *EventData) {
	for _, task := range workflow.listeners {
		in := task.convert(event)
		d, err := json.Marshal(in)
		if err != nil {
			log.Fatalln(err)
		}
		_, err = cli.ExecuteTask(context.Background(), &ExecuteTaskRequest{
			Service:  task.service,
			TaskData: string(d),
			TaskKey:  task.task,
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
}
