package main

import (
	"github.com/mesg-foundation/engine/protobuf/api"
	pb "github.com/mesg-foundation/engine/protobuf/api"
	"github.com/mesg-foundation/engine/service"
)

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

func newTestCreateServiceRequest() *pb.CreateServiceRequest {
	return &api.CreateServiceRequest{
		Sid:  "test-service",
		Name: "test-service",
		Configuration: service.Service_Configuration{
			Env: []string{"FOO=1", "BAR=2", "REQUIRED"},
		},
		Tasks: []*service.Service_Task{
			{
				Key: "task1",
				Inputs: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "String",
					},
				},
				Outputs: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "String",
					},
					{
						Key:  "timestamp",
						Type: "Number",
					},
				},
			},
			{
				Key: "task2",
				Inputs: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "String",
					},
				},
				Outputs: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "String",
					},
					{
						Key:  "timestamp",
						Type: "Number",
					},
				},
			},
			{
				Key: "task_complex",
				Inputs: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "Object",
						Object: []*service.Service_Parameter{
							{
								Key:  "msg",
								Type: "String",
							},
							{
								Key:      "array",
								Type:     "String",
								Repeated: true,
								Optional: true,
							},
						},
					},
				},
				Outputs: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "Object",
						Object: []*service.Service_Parameter{
							{
								Key:  "msg",
								Type: "String",
							},
							{
								Key:  "timestamp",
								Type: "Number",
							},
							{
								Key:      "array",
								Type:     "String",
								Repeated: true,
								Optional: true,
							},
						},
					},
				},
			},
		},
		Events: []*service.Service_Event{
			{
				Key: "test_service_ready",
			},
			{
				Key: "test_event",
				Data: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "String",
					},
					{
						Key:  "timestamp",
						Type: "Number",
					},
				},
			},
			{
				Key: "test_event_complex",
				Data: []*service.Service_Parameter{
					{
						Key:  "msg",
						Type: "Object",
						Object: []*service.Service_Parameter{
							{
								Key:  "msg",
								Type: "String",
							},
							{
								Key:  "timestamp",
								Type: "Number",
							},
							{
								Key:      "array",
								Type:     "String",
								Repeated: true,
							},
						},
					},
				},
			},
			{
				Key: "event_after_task",
				Data: []*service.Service_Parameter{
					{
						Key:  "task_key",
						Type: "String",
					},
					{
						Key:  "timestamp",
						Type: "Number",
					},
				},
			},
		},
		Source: "QmWHKNvJ4wT83TLHPLMjTRBjJYvcwyr8oqTtqxaJZVXPbQ",
	}
}
