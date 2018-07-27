# Listen for events

## Why listen for events?

Applications listen for events in real time to execute actions or tasks when something relevant happens on a technology.

## Listening for events from services

To listen for events, the application needs to open a stream with Core with [gRPC](https://grpc.io/) using the [Protobuffer definition](https://github.com/mesg-foundation/core/blob/dev/api/core/api.proto). When opening the stream, the application listens to the service. It can listen to many service as the same time.

<tabs>
<tab title="Request" vp-markdown>

### `Core.ListenEvent`

| **Name** | **Type** | **Required** | **Description** |
| --- | --- | --- | --- |
| **serviceID** | `String` | Required | ID of the service that you want to listen. |
| **eventFilter** | `String` | Optional | Only listens for this given event ID. |

```javascript
{
  "serviceID": "v1_fe25be776e1e256400c77067a1cb7666",
  "eventFilter": "eventIDToOnlyListenTo"
}
```

</tab>

<tab title="Stream Reply" vp-markdown>

| **Name** | **Type** | **Description** |
| --- | --- | --- | --- |
| **eventKey** | `String` | The event's key defined in the [service file](../service/service-file.md). |
| **eventData** | `String` | The event's data in JSON format. |

```javascript
{
  "eventKey": "eventX",
  "eventData": "{\"dataX\": \"event data\"}"
}
```

</tab>
</tabs>

### Example

<tabs>
<tab title="Node" vp-markdown>

```javascript
const MESG = require('mesg-js').application()

MESG.api.ListenEvent({
  serviceID: "v1_fe25be776e1e256400c77067a1cb7666",
})
.on('error', function(error) {
  // An error has occurred and the stream has been closed.
})
.on('data', function(data) {
  ...
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
    stream, _ := cli.ListenEvent(context.Background(), &core.ListenEventRequest{
        ServiceID: "v1_fe25be776e1e256400c77067a1cb7666",
    })
    for {
        event, _ := stream.Recv()
        fmt.Println(event.EventKey, event.EventData)
        // TODO process event
    }
}
```

</tab>
</tabs>

## Listen for task execution outputs

The execution of tasks can take a long time to finish depending on the action they are completing, so outputs are sent back asynchronously. To listen for task execution's outputs, applications need to open a stream with the Core through [gRPC](https://grpc.io/) and use the [Protobuffer definition](https://github.com/mesg-foundation/core/blob/dev/api/core/api.proto).

::: warning
Outputs are sent asynchronously. Make sure that the application listens for outputs before it executes a task, otherwise it will miss the outputs.
:::

<tabs>
<tab title="Request" vp-markdown>

### `Core.ListenResult`

| **Name** | **Type** | **Required** | **Description** |
| --- | --- | --- | --- |
| **serviceID** | `String` | Required | ID of the service. |
| **taskFilter** | `String` | Optional | Only listens for this given task ID. |
| **outputFilter** | `String` | Optional | Only listens for this given output ID. If set, the attribute `taskFilter` should also be provided. |

```javascript
{
  "serviceID": "v1_fe25be776e1e256400c77067a1cb7666",
  "taskFilter": "taskIDToOnlyListenTo",
  "outputFilter": "outputIDToOnlyListenTo"
}
```

</tab>

<tab title="Stream Reply" vp-markdown>

| **Name** | **Type** | **Description** |
| --- | --- | --- | --- | --- | --- |
| **executionID** | `String` | The execution ID of this output. |
| **taskKey** | `String` | The key of the task as defined in the [service file](../service/service-file.md). |
| **outputKey** | `String` | The key of the output of the task as defined in the [service file](../service/service-file.md). |
| **outputData** | `String` | The data returned by the task serialized in JSON. |

```javascript
{
  "executionID": "xxxxx",
  "taskKey": "taskX",
  "outputKey": "outputX",
  "outputData": "{\"outputValX\": \"result of execution\"}",
}
```

</tab>
</tabs>

### Examples

<tabs>
<tab title="Node" vp-markdown>

```javascript
const MESG = require('mesg-js').application()

MESG.api.ListenResult({
  serviceID: "v1_31341hj33189yr13iugr13",
})
.on('error', function(error) {
  // An error has occurred and the stream has been closed.
})
.on('data', function(data) {
  ...
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
    stream, _ := cli.ListenResult(context.Background(), &core.ListenResultRequest{
        ServiceID: "v1_fe25be776e1e256400c77067a1cb7666",
    })
    for {
        result, _ := stream.Recv()
        fmt.Println(result.ExecutionID, result.OutputKey, result.OutputData)
    }
}
```

</tab>
</tabs>

