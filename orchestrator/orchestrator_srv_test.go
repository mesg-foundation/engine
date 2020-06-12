package orchestrator

import (
	"context"

	sdktypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/mesg-foundation/engine/hash"
	"github.com/mesg-foundation/engine/service"
)

func createTestService(store *storeTest) (hash.Hash, error) {
	taskPrice := &service.Service_Task_Price{
		PerCall: sdktypes.NewInt(0),
		PerKB:   sdktypes.NewInt(0),
		PerSec:  sdktypes.NewInt(0),
	}
	task1Price := &service.Service_Task_Price{
		PerCall: sdktypes.NewInt(1000),
		PerKB:   sdktypes.NewInt(1000),
		PerSec:  sdktypes.NewInt(30000),
	}
	return store.CreateService(context.Background(),
		"test-service",
		"test-service",
		"",
		service.Service_Configuration{
			Env: []string{"FOO=1", "BAR=2", "REQUIRED"},
		}, []*service.Service_Task{
			{
				Key:   "task_trigger",
				Price: taskPrice,
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
				Key:   "task1",
				Price: task1Price,
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
				Key:   "task2",
				Price: taskPrice,
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
				Key:   "task_complex",
				Price: taskPrice,
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
				Key:   "task_complex_trigger",
				Price: taskPrice,
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
			}}, []*service.Service_Event{
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
			}},
		nil,
		"",
		"NOT_NEEDED",
	)
}
