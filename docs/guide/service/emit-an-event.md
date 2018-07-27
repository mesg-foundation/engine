# Emit an Event

## Why emit an Event?

Events are emitted from a Service \(e.g.: a web server receiving a request, or a blockchain technology receiving a new transaction\). These events are emitted to achieve a desired effect or to be used as a trigger to make another task happen. Each Service has different kinds of events that you can send to Core.

## Steps to follow

To emit events from your Service, you'll need to:

* [Add the definition of the event](#event-definitions) in the Service's [`mesg.yml`](service-file.md) file
* [Emit the event](#emit-an-event-2) when it happens on the Service

## Event definitions

To create an event, the first step is to update the Service's [`mesg.yml`](service-file.md) file and add an event indexed by its key with the following attributes:

| **Attribute** | **Default value** | **Type** | **Description** |
| --- | --- | --- | --- |
| **name** | `id` | `String` | Name of the event. If the name is not set, it will be the same as the ID you choose for the event. |
| **description** | `""` | `String` | Describe the event, what's its purpose is and why users would want to use it. |
| **data** | `{}` | `map<id,`[`Data`](emit-an-event.md#data-of-your-event)`>` | The structure of the event's data. |

### Event's data

| **Attribute** | **Default value** | **Type** | **Description** |
| --- | --- | --- | --- | --- |
| **name** | `id` | `String` | Name of the data |
| **description** | `""` | `String` | Description of the data |
| **type** | `String` | [`Type`](emit-an-event.md#type-of-your-data) | Type of data |
| **optional** | `false` | `boolean` | Mark the data as optional |

### Data's type

The data type can be one of the following:

* `String`
* `Boolean`
* `Number`
* `Object`

### Example

Example of an event definition in a [`mesg.yml`](service-file.md) file:

```yaml
...
events:
    eventX:
        name: "Event X"
        description: "This is event X"
        data:
            foo:
                name: "Foo"
                description: "Foo is the first data"
                type: String
                optional: false
            bar:
                name: "Bar"
                description: "Bar is the second data"
                type: Boolean
                optional: true
...
```

## Emit an Event

To emit events from the Service to the Core, the Service has to follow the [Protobuffer definition](https://github.com/mesg-foundation/core/blob/dev/api/service/api.proto) and use [gRPC](https://grpc.io/).

::: tip
Consider emitting event when the service is ready. If the service needs to synchronize data first, it should wait for the synchronization to complete before emitting events.
:::

<tabs>
<tab title="Request" vp-markdown>

### `Service.EmitEvent`

| **Name** | **Type** | **Required** | **Description** |
| --- | --- | --- | --- |
| **token** | `String` | Required | The token given by the Core as environment variable `MESG_TOKEN` |
| **eventKey** | `String` | Required | The event's key defined in the [service file](https://github.com/mesg-foundation/documentation/tree/c1028b6f9d709adf2ad46364ce7baaa37e27ff8e/service/service/service-file.md) |
| **eventData** | `String` | Required | The event's data in JSON format |

```javascript
{
    "token": "TOKEN_FROM_ENV",
    "eventKey": "eventX",
    "eventData": "{\"foo\":\"hello\",\"bar\":false}"
}
```

</tab>

<tab title="Reply" vp-markdown>

| **Name** | **Type** | **Description** |
| --- | --- | --- |


```javascript
{}
```

</tab>
</tabs>

### Examples

<tabs>
<tab title="Node" vp-markdown>

```javascript
const MESG = require('mesg-js').service()

MESG.emitEvent("eventX", {
  foo: "hello",
  bar: false,
})
```

</tab>

<tab title="Go" vp-markdown>

```go
package main

import (
    "context"
    "encoding/json"
    "io/ioutil"
    "log"
    "os"

    api "github.com/mesg-foundation/core/api/service"
    "google.golang.org/grpc"
    yaml "gopkg.in/yaml.v2"
)

type EventX struct {
    Foo string
    Bar bool
}

func main() {
    connection, _ := grpc.Dial(os.Getenv("MESG_ENDPOINT"), grpc.WithInsecure())
    cli := api.NewServiceClient(connection)

    eventX, _ := json.Marshal(EventX{
        Foo: "hello",
        Bar: false,
    })

    reply, _ := cli.EmitEvent(context.Background(), &api.EmitEventRequest{
        Token:   os.Getenv("MESG_TOKEN"),
        EventKey:  "eventX",
        EventData: string(eventX),
    })
    log.Println(reply)
}
```

</tab>
</tabs>

