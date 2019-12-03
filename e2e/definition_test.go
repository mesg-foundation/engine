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
		Source: "QmZ8r6Basr67m1FYe9mar3B9gB1kDfDZnqZvs3EVj46Np4",
	}
}

func newTestComplexCreateServiceRequest() *pb.CreateServiceRequest {
	return &api.CreateServiceRequest{
		Sid:  "test-complex-service",
		Name: "test-complex-service",
		Dependencies: []*service.Service_Dependency{
			{
				Key:     "nginx",
				Image:   "nginx",
				Volumes: []string{"/etc/nginx"},
			},
		},
		Configuration: service.Service_Configuration{
			Env: []string{
				"ENVA=do_not_override",
				"ENVB=override",
			},
			Volumes:     []string{"/volume/test/"},
			VolumesFrom: []string{"nginx"},
		},
		Events: []*service.Service_Event{
			{Key: "test_service_ready"},
			{Key: "read_env_ok"},
			{Key: "read_env_error"},
			{Key: "read_env_override_ok"},
			{Key: "read_env_override_error"},
			{Key: "access_volumes_ok"},
			{Key: "access_volumes_error"},
			{Key: "access_volumes_from_ok"},
			{Key: "access_volumes_from_error"},
			{Key: "resolve_dependence_ok"},
			{Key: "resolve_dependence_error"},
		},
		Source: "QmSuVcdic2dhS5QKQGWp66SJQUkDRqAqCHpU6Sx9uXJcdc",
	}
}
