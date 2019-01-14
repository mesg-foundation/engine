# Create your application

Now that you have your different Services ready and deployed, you need to connect them through an Application. You will be able to connect a task to an event with libraries that we provide, otherwise you can always directly connect to MESG Core to [listen for events](listen-for-events.md) and [execute tasks](execute-a-task.md).

<tabs>
<tab title="Node" vp-markdown>

```javascript
const MESG = require('mesg-js').application()

// When SERVICE_EVENT_ID emits event "eventX"
// then execute "taskX" from SERVICE_TASK_ID 
MESG.listenEvent({ serviceID: SERVICE_EVENT_ID, eventFilter: 'eventX' })
  .on('data', (event) => {
    MESG.executeTask({
      serviceID: SERVICE_TASK_ID,
      taskKey: 'taskX',
      inputData: JSON.stringify({ foo: 'bar' })
    }).catch((err) => console.log(err.message))
  })

// When SERVICE_TASK_ID send the result of taskX
// then execute "taskB" from SERVICE_TASK2_ID
MESG.listenResult({ serviceID: SERVICE_TASK_ID, taskFilter: 'taskX' })
  .on('data', (event) => {
    MESG.executeTask({
      serviceID: SERVICE_TASK2_ID,
      taskKey: 'taskB',
      inputData: JSON.stringify({ hello: "world" })
    }).catch((err) => console.log(err.message))
  })
```

[See the MESG.js library for additional documentation](https://github.com/mesg-foundation/mesg-js/tree/master#application)

</tab>

<tab title="Go" vp-markdown>

```go
package main

import (
	"context"
	"encoding/json"

	"github.com/mesg-foundation/core/protobuf/coreapi"
	"google.golang.org/grpc"
)

func main() {
	connection, _ := grpc.Dial(":50052", grpc.WithInsecure())
	client := coreapi.NewCoreClient(connection)

	go func() {
		stream, _ := client.ListenEvent(context.Background(), &coreapi.ListenEventRequest{
			ServiceID:   SERVICE_EVENT_ID,
			EventFilter: "eventX",
		})

		for {
			event, _ := stream.Recv()
			inputData, _ := json.Marshal(map[string]string{"foo": "bar"})
			client.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
				ServiceID: SERVICE_TASK_ID,
				TaskKey:   "taskX",
				InputData: string(inputData),
			})
		}
	}()

	go func() {
		stream, _ := client.ListenResult(context.Background(), &coreapi.ListenEventRequest{
			ServiceID: SERVICE_TASK_ID,
			TaskKey:   "taskX",
		})

		for {
			event, _ := stream.Recv()
			inputData, _ := json.Marshal(map[string]string{"hello": "world"})
			client.ExecuteTask(context.Background(), &coreapi.ExecuteTaskRequest{
				ServiceID: SERVICE_TASK2_ID,
				TaskKey:   "taskB",
				InputData: string(inputData),
			})
		}
	}()
}

```

</tab>
</tabs>

These kinds of workflows should be enough for most use cases, but if you want to create more complex applications, you can connect directly to APIs and keep reading the documentation.

::: tip Get Help
You need help ? Check out the <a href="https://forum.mesg.com" target="_blank">MESG Forum</a>.
