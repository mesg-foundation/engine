




# Service API
<!--
DO NOT EDIT
This file is generated using the ./scripts/build-proto.sh scripts
Please update the Service file
-->

This is the API for MESG Services to interact with MESG Core.
It is to be consumed only by MESG Services.
It provides all necessary functions that MESG Services need in order to interact with MESG Core.

This API is only accessible through [gRPC](https://grpc.io/).

Applications must not use this API, but rather use the [Core API](./core.md).

The source file of this API is hosted on [GitHub](https://github.com/mesg-foundation/core/blob/master/api/service/api.proto).

[[toc]]


## EmitEvent

Emit an event to Core.
The event and its data must be defined in the [service's definition file](../guide/service/service-file.md).

<tabs>
<tab title="Request">





#### EmitEventRequest
The request's data for the `EmitEvent` API.

**Example:**
```json
{
  "token":     "__SERVICE_TOKEN_FROM_ENV__",
  "eventKey":  "__EVENT_KEY__",
  "eventData": "{\"foo\":\"hello\",\"bar\":false}"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| token | [string](#string) | The service's token given by Core as the environment variable `MESG_TOKEN`. |
| eventKey | [string](#string) | The event's key as defined in the [service file](../guide/service/service-file.md). |
| eventData | [string](#string) | The event's data encoded in JSON as defined in the [service file](../guide/service/service-file.md). |













</tab>

<tab title="Reply">



#### EmitEventReply
Reply of `EmitEvent` API doesn't contain any data.















</tab>
</tabs>

## ListenTask

Subscribe to the stream of tasks to execute.
Every task received must be executed and its result must be submitted using the `SubmitResult` API.

<tabs>
<tab title="Request">







#### ListenTaskRequest
The request's data for the `ListenTask` stream API.

**Example:**
```json
{
  "token": "__SERVICE_TOKEN_FROM_ENV__"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| token | [string](#string) | The service's token given by the Core as the environment variable `MESG_TOKEN`. |











</tab>

<tab title="Reply">













#### TaskData
The data received from the stream of the `ListenTask` API.
The data will be received over time as long as the stream is open.
The `executionID` value must be kept and sent with the result when calling the [`SubmitResult` API](#submitresult).

**Example:**
```json
{
  "executionID": "__EXECUTION_ID__",
  "taskKey":     "__TASK_KEY__",
  "inputData":   "{\"inputX\":\"Hello world!\",\"inputY\":true}"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| executionID | [string](#string) | The unique identifier of the execution. Must be kept and sent with the result when calling the [`SubmitResult` API](#submitresult). |
| taskKey | [string](#string) | The key from the task to execute as defined in the [service file](../guide/service/service-file.md). |
| inputData | [string](#string) | The task's input encoded in JSON as defined in the [service file](../guide/service/service-file.md). |





</tab>
</tabs>

## SubmitResult

Submit the result of a task's execution to Core.
The result must be defined as a task's output in the [service's definition file](../guide/service/service-file.md).

<tabs>
<tab title="Request">











#### SubmitResultRequest
The request's data for the `SubmitResult` API.
The data must contain the `executionID` of the executed task received from the stream of [`ListenTask` API](#listentask).

**Example:**
```json
{
  "executionID": "__EXECUTION_ID__",
  "outputKey":   "__OUTPUT_KEY__",
  "outputData":  "{\"foo\":\"super result\",\"bar\":true}"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| executionID | [string](#string) | The `executionID` received from the [`ListenTask` stream](#listentask). |
| outputKey | [string](#string) | The output's key as defined in the [service file](../guide/service/service-file.md). |
| outputData | [string](#string) | The result's data encoded in JSON as defined in the [service file](../guide/service/service-file.md). |







</tab>

<tab title="Reply">









#### SubmitResultReply
Reply of `SubmitResult` API doesn't contain any data.









</tab>
</tabs>



