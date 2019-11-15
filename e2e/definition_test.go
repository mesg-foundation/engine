package main

import (
	"github.com/mesg-foundation/engine/protobuf/api"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/service"
)

func newTestCreateServiceRequest() *pb.CreateServiceRequest {
	return &api.CreateServiceRequest{
		Sid:  "test-service",
		Name: "test-service",
		Configuration: service.Service_Configuration{
			Env: []string{"FOO=1", "BAR=2", "REQUIRED"},
		},
		Tasks: []*service.Service_Task{
			{
				Key: "ping",
				Inputs: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "String",
					},
				},
				Outputs: []*service.Service_Parameter{
					{
						Key:  "pong",
						Type: "String",
					},
				},
			},
			{
				Key: "add",
				Inputs: []*service.Service_Parameter{
					{
						Key:  "n",
						Type: "Number",
					},
					{
						Key:  "m",
						Type: "Number",
					},
				},
				Outputs: []*service.Service_Parameter{
					{
						Key:  "res",
						Type: "Number",
					},
				},
			},
			{
				Key: "error",
			},
		},
		Events: []*service.Service_Event{
			{
				Key: "test_service_ready",
			},
			{
				Key: "ping_ok",
				Data: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "String",
					},
				},
			},
			{
				Key: "add_ok",
				Data: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "String",
					},
				},
			},
			{
				Key: "error_ok",
				Data: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "String",
					},
				},
			},
		},
		Source: "QmPG1Ze96pH1EgVMWsGKM33jXoG63rigMncSEqZXP7oncq",
	}
}

func newTestComplexCreateServiceRequest() *pb.CreateServiceRequest {
	return &api.CreateServiceRequest{
		Sid:  "test-complex-service",
		Name: "test-complex-service",
		Dependencies: []*service.Service_Dependency{
			{
				Key:     "busybox",
				Image:   "busybox",
				Volumes: []string{"/tmp"},
			},
		},
		Configuration: service.Service_Configuration{
			VolumesFrom: []string{"busybox"},
			Env:         []string{"FOO"},
		},
		Tasks: []*service.Service_Task{
			{
				Key: "foo",
				Inputs: []*service.Service_Parameter{
					{
						Key:  "o",
						Type: "Object",
						Object: []*service.Service_Parameter{
							{
								Key:  "msg",
								Type: "String",
							},
						},
					},
				},
			},
		},
		Events: []*service.Service_Event{
			{
				Key: "foo",
				Data: []*service.Service_Parameter{
					{
						Key:  "o",
						Type: "Object",
						Object: []*service.Service_Parameter{
							{
								Key:  "msg",
								Type: "String",
							},
						},
					},
				},
			},
		},
		Source: "QmRarLspkiAUL7nQiTfA2DGTusVpY7JUB33Qm9poq56ris",
	}
}
