package main

import (
	"github.com/mesg-foundation/engine/service"
	servicemodule "github.com/mesg-foundation/engine/x/service"
)

var testComplexCreateServiceMsg = &servicemodule.MsgCreate{
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
		{Key: "service_ready"},
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
	Source: "QmQbobdv3tg1UdPiHBGaqbQF2mjGHPBBNAK4M9hoFA3NVU",
}

var testCreateServiceMsg = &servicemodule.MsgCreate{
	Sid:  "test-service",
	Name: "test-service",
	Configuration: service.Service_Configuration{
		Env: []string{"FOO=1", "BAR=2", "REQUIRED"},
	},
	Tasks: []*service.Service_Task{
		{
			Key: "task_trigger",
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
		{
			Key: "task_complex_trigger",
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
			Key: "service_ready",
		},
		{
			Key: "event_trigger",
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
			Key: "event_complex_trigger",
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
	},
	Source: "QmXiBMy2rguhPn5hCmHh8rw1jjpq6YcekiUcvh91htVndM",
}
