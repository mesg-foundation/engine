# Execute a task

## Why execute a service's task?

Applications can execute a service's task allowing you to reuse the maximum number of already-built logic and enjoy the MESG ecosystem.

## Execute a service's task

To execute a task, applications need to connect to Core through [gRPC](https://grpc.io/) and use the [Protobuffer definition](https://github.com/mesg-foundation/core/blob/dev/api/core/api.proto). Core will reply with an `executionID`that identifies the task's execution. To get the output of the task's execution, the application has to listen for an [execution output.](execute-a-task.md#listen-for-execution-outputs)

<tabs>
<tab title="Request" vp-markdown>

### `Core.ExecuteTask`

| **Name** | **Type** | **Required** | **Description** |
| --- | --- | --- | --- |
| **serviceId** | `String` | Required | ID of the service. |
| **taskKey** | `String` | Required | The task's key defined in the [service file](../service/service-file.md). |
| **inputData** | `String` | Required | The task's inputs in JSON format. |

```javascript
{
  "serviceID": "v1_fe25be776e1e256400c77067a1cb7666",
  "taskKey": "taskX",
  "inputData": "{\"inputX\":\"input value\"}"
}
```

</tab>

<tab title="Reply" vp-markdown>

| **Name** | **Type** | **Description** |
| --- | --- | --- |
| **executionID** | `String` | The ID of the execution. |

```javascript
{
  "executionID": "xxxxx"
}
```

</tab>
</tabs>

### Examples

<tabs>
<tab title="Node" vp-markdown>

```javascript
const MESG = require('mesg-js').application()

MESG.api.ExecuteTask({
  serviceID: "v1_fe25be776e1e256400c77067a1cb7666",
  taskKey: "taskX",
  inputData: JSON.stringify({
    inputX: "input value"
  })
}, (err, reply) => {
  // handle response if needed
})
```

</tab>

<tab title="Go" vp-markdown>

```go
package main

import (
    "context"
    "fmt"

    "github.com/mesg-foundation/core/api/core"
    "github.com/mesg-foundation/core/service"
    "google.golang.org/grpc"
)

func main() {
    connection, _ := grpc.Dial(":50052", grpc.WithInsecure())
    cli := core.NewCoreClient(connection)
    res, _ := cli.ExecuteTask(context.Background(), &core.ExecuteTaskRequest{
        ServiceID:  "v1_fe25be776e1e256400c77067a1cb7666",
        TaskKey:  "taskX",
        InputData: "{\"inputX\":\"input value\"}",
    })
    fmt.Println(res.ExecutionID)
}
```

</tab>
</tabs>

::: tip Get Help
You need help ? Check out the <a href="https://forum.mesg.com" target="_blank">MESG Forum</a>.