




# Core API
<!--
DO NOT EDIT
This file is generated using the ./scripts/build-proto.sh scripts
Please update the Core file
-->

This is the primary API to interact with MESG Core functionalities.
It can be consumed by any applications or tools that you'd like to interact with MESG Core.
It is actually used by the MESG CLI and MESG Application libraries.

This API is only accessible through [gRPC](https://grpc.io/).

Services must not use this API, but rather use the [Service API](./service.md).

The source file of this API is hosted on [GitHub](https://github.com/mesg-foundation/core/blob/master/api/core/api.proto).

[[toc]]


## ListenEvent

Subscribe to a stream that listens for events from a service.

<tabs>
<tab title="Request">































#### ListenEventRequest
The request's data for the `ListenEvent` stream's API.

**Example**
```json
{
  "serviceID":   "__SERVICE_ID__",
  "eventFilter": "__EVENT_KEY_TO_MATCH__"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| serviceID | [string](#string) | The Service ID. Generated when using the [`DeployService` API](#deployservice). |
| eventFilter | [string](#string) | __Optional.__ Event's key to filter. The event must match this key. The default is `*` which matches any event. |





























</tab>

<tab title="Reply">

















#### EventData
The data received from the stream of the `ListenEvent` API.
The data will be received over time as long as the stream is open.

**Example**
```json
{
  "eventKey":  "__EVENT_KEY__",
  "eventData": "{\"foo\":\"bar\"}"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| eventKey | [string](#string) | The event's key. |
| eventData | [string](#string) | The event's data encoded in JSON. |











































</tab>
</tabs>

## ListenResult

Subscribe to a stream that listens for task's result from a service.

<tabs>
<tab title="Request">

































#### ListenResultRequest
The request's data for the `ListenResult` stream API.

**Example**
```json
{
  "serviceID":     "__SERVICE_ID__",
  "taskFilter":    "__TASK_KEY_TO_MATCH__",
  "outputFilter":  "__OUTPUT_KEY_TO_MATCH__",
  "tagFilters":     ["tagX"]
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| serviceID | [string](#string) | The Service ID. Generated when using the [`DeployService` API](#deployservice). |
| taskFilter | [string](#string) | __Optional.__  The task's key to filter. The task must match this key. The default is `*` which matches any task. |
| outputFilter | [string](#string) | __Optional.__ The output's key from the task to filter. The task must return this output's key. The default is `*` which matches any output. |
| tagFilters | [string](#string)[] | __Optional.__ The list of tags to filter. This is a "match all" list. All tags in parameters should be included in the execution to match. |



























</tab>

<tab title="Reply">









































#### ResultData
The data received from the stream of the `ListenResult` API.
The data will be received over time as long as the stream is open.

**Example**
```json
{
  "executionID":   "__EXECUTION_ID__",
  "taskKey":       "__TASK_KEY__",
  "outputKey":     "__OUTPUT_KEY__",
  "outputData":    "{\"foo\":\"bar\"}",
  "executionTags": ["executionX", "test"]
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| executionID | [string](#string) | The unique identifier of the execution. |
| taskKey | [string](#string) | The key of the executed task. |
| outputKey | [string](#string) | The output's key from the returned task. |
| outputData | [string](#string) | The output's data from the returned task, encoded in JSON. |
| executionTags | [string](#string)[] | The list of tags associated with the execution |



















</tab>
</tabs>

## ExecuteTask

Execute a service's task through [Core](../guide/start-here/core.md).

<tabs>
<tab title="Request">





















#### ExecuteTaskRequest
The request's data for the `ExecuteTask` API.

**Example**
```json
{
  "serviceID":     "__SERVICE_ID__",
  "taskKey":       "__TASK_KEY__",
  "inputData":     "{\"foo\":\"bar\"}",
  "executionTags": ["executionX", "test"]
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| serviceID | [string](#string) | The Service ID. Generated when using the [`DeployService` API](#deployservice). |
| taskKey | [string](#string) | The task's key to execute. |
| inputData | [string](#string) | The inputs of the task to execute, encoded in JSON. |
| executionTags | [string](#string)[] | __Optional.__ The list of tags to associate with the execution |







































</tab>

<tab title="Reply">



















#### ExecuteTaskReply
The reply's data of the `ExecuteTask` API.

**Example**
```json
{
  "executionID": "__EXECUTION_ID__"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| executionID | [string](#string) | The unique identifier of the execution. |









































</tab>
</tabs>

## StartService

Start a service. The service must be already deployed to [Core](../guide/start-here/core.md).

<tabs>
<tab title="Request">

















































#### StartServiceRequest
The request's data for the `StartService` API.

**Example**
```json
{
  "serviceID": "__SERVICE_ID__"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| serviceID | [string](#string) | The Service ID. Generated when using the [`DeployService` API](#deployservice). |











</tab>

<tab title="Reply">















































#### StartServiceReply
Reply of `StartService` API doesn't contain any data.













</tab>
</tabs>

## StopService

Stop a service. The service must be already deployed to [Core](../guide/start-here/core.md).

<tabs>
<tab title="Request">





















































#### StopServiceRequest
The request's data for the `StopService` API.

**Example**
```json
{
  "serviceID": "__SERVICE_ID__"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| serviceID | [string](#string) | The Service ID. Generated when using the [`DeployService` API](#deployservice). |







</tab>

<tab title="Reply">



















































#### StopServiceReply
Reply of `StopService` API doesn't contain any data.









</tab>
</tabs>

## DeployService

Deploy a service to [Core](../guide/start-here/core.md). This will give you an unique identifier which is used to interact with the service.

<tabs>
<tab title="Request">













#### DeployServiceRequest
The data sent to the request stream of the `DeployService` API.
Stream should be closed after url or all chunks sent to server.

**Example**
```json
{
  "url": "__SERVICE_GIT_URL__"
}
```
or
```json
{
  "chunk": __SERVICE_GZIPPED_TAR_FILE_CHUNK__
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| url | [string](#string) | Git repo url of service. When url provided, stream will be closed after the first receive. |
| chunk | [bytes](#bytes) | Chunks of gzipped tar archive of service. If chunk provided, stream should be closed by client after all chunks sent. |















































</tab>

<tab title="Reply">









#### DeployServiceReply
The data received from the reply stream of the `DeployService` API.
Stream will be closed by server after deployment is done.

**Example**
```json
{
  "status": {
    message: "__STATUS_MESSAGE__",
    type: __STATUS_TYPE__,
  }
}
```
or
```json
{
  "serviceID": "__SERVICE_ID__"
}
```
or
```json
{
  "validationError": "__SERVICE_VALIDATION_ERROR__"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| status | [DeployServiceReply.Status](#api.DeployServiceReply.Status) | `status` will be sent after each status change. |
| serviceID | [string](#string) | `serviceID` will be sent as the last message of stream when service deployed successfully. |
| validationError | [string](#string) | `validationError` will be sent as the last message of stream when there is a validation error. |



















































</tab>
</tabs>

## DeleteService

Delete a service from Core. This function only deletes a deployed service in [Core](../guide/start-here/core.md). If the service's code is on your computer, the source code will not be deleted.

<tabs>
<tab title="Request">





#### DeleteServiceRequest
Request's data of the `DeleteService` API.

**Example**
```json
{
  "serviceID": "__SERVICE_ID__"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| serviceID | [string](#string) | The Service ID. Generated when using the [`DeployService` API](#deployservice). |























































</tab>

<tab title="Reply">



#### DeleteServiceReply
Reply of `DeleteService` API doesn't contain any data.

























































</tab>
</tabs>

## ListServices

List all services already deployed in [Core](../guide/start-here/core.md).

<tabs>
<tab title="Request">





























#### ListServicesRequest
Reply of `ListServices` API doesn't contain any data.































</tab>

<tab title="Reply">



























#### ListServicesReply
The reply's data of the `ListServices` API.

**Example**
```json
[{
  "service": {
    "id": "idX",
    "name": "serviceX",
    "description: "descriptionX",
    "status: "statusX",
    "events": {
      "eventX": {
        "data": {
          "dataX": { "type": "String" }
        }
      }
    },
    "tasks": {
      "taskX": {
        "inputs": {
          "foo": { "type": "String" }
        },
        "outputs": {
          "outputX": {
            "data": {
              "resX": { "type": "String" }
            }
          }
        }
      }
    }
  }
}]
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| services | [Service](#api.Service)[] | The list of previously-deployed services' definitions. |

































</tab>
</tabs>

## GetService

Get the definition of an already-deployed service from its ID.

<tabs>
<tab title="Request">

























#### GetServiceRequest
The request's data for the `GetService` API.

**Example**
```json
{
  "serviceID": "__SERVICE_ID__"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| serviceID | [string](#string) | The Service ID. Generated when using the [`DeployService` API](#deployservice). |



































</tab>

<tab title="Reply">























#### GetServiceReply
The reply's data of the `GetService` API.

**Example**
```json
{
  "service": {
    "id": "idX",
    "name": "serviceX",
    "description: "descriptionX",
    "status: "statusX",
    "events": {
      "eventX": {
        "data": {
          "dataX": { "type": "String" }
        }
      }
    },
    "tasks": {
      "taskX": {
        "inputs": {
          "foo": { "type": "String" }
        },
        "outputs": {
          "outputX": {
            "data": {
              "resX": { "type": "String" }
            }
          }
        }
      }
    }
  }
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| service | [Service](#api.Service) | Service's definition. |





































</tab>
</tabs>

## ServiceLogs

ServiceLogs gives a stream for dependency logs of a service.

<tabs>
<tab title="Request">













































#### ServiceLogsRequest
The request's data for `ServiceLogs` API.

**Example**
```json
{
  "serviceID": "__SERVICE_ID__",
  "dependencies": "__SERVICE_DEPENDENCIES__"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| serviceID | [string](#string) | The Service ID. Generated when using the [`DeployService` API](#deployservice). |
| dependencies | [string](#string)[] | __Optional.__ List of dependencies to filter service logs. All by default. |















</tab>

<tab title="Reply">



































#### LogData
The data received from the stream of the `ServiceLogs` API.
The data will be received over time as long as the stream is open.

**Example**
```json
{
  "dependency":  "__SERVICE_DEPENDENCY__",
  "type": __LOG_TYPE__,
  "data":  __LOG_CHUNK__,
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| dependency | [string](#string) | Service dependency that data belongs. |
| type | [LogData.Type](#api.LogData.Type) | The log type. |
| data | [bytes](#bytes) | Log data chunk. |

























</tab>
</tabs>



