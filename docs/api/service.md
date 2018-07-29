


# Service API
<!--
DO NOT EDIT
This file is generated using the ./scripts/build-proto.sh scripts
Please update the Service file
-->

## APIs

API accessible from services.
This API can and should only be accessible when you create a MESG Service.
It provide all functions necessary to be able to execute tasks and submit specific events.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| EmitEvent | [EmitEventRequest](#EmitEventRequest) | [EmitEventReply](#EmitEventRequest) | Let you to emit an event to the [Core](/guide/start-here/core.html) based on the ones defined in your [mesg.yml](/guide/service/service-file.html). |
| ListenTask | [ListenTaskRequest](#ListenTaskRequest) | [TaskData](#ListenTaskRequest) | Subscribe to the stream of tasks that will receive tasks from the [Core](/guide/start-here/core.html) |
| SubmitResult | [SubmitResultRequest](#SubmitResultRequest) | [SubmitResultReply](#SubmitResultRequest) | Let you submit a result from a task to the [Core](/guide/start-here/core.html). The result should be an output of the tasks |




## Types



<a name="EmitEventReply"/>

### EmitEventReply
Response of the Core when receiving an event from the `EmitEvent` call

**Example:**
```json
{}
```




<a name="EmitEventRequest"/>

### EmitEventRequest
Data sent while calling the `EmitEvent` API.

**Example:**
```json
{
  "token": "TOKEN_FROM_ENV",
  "eventKey": "eventX",
  "eventData": "{\"foo\":\"hello\",\"bar\":false}"
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token | [string](#string) |  | The token given by the Core as environment variable `MESG_TOKEN`. |
| eventKey | [string](#string) |  | The event's key defined in the [service file](/guide/service/service-file.html). |
| eventData | [string](#string) |  | The data of your event encoded in JSON as defined in your [mesg.yml](/guide/service/service-file.html). |




<a name="ListenTaskRequest"/>

### ListenTaskRequest
Data sent to connect to the `ListenTask` stream API.

**Example:**
```json
{
  "token": "TOKEN_FROM_ENV"
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token | [string](#string) |  | The token given by the Core as environment variable `MESG_TOKEN`.` |




<a name="SubmitResultReply"/>

### SubmitResultReply
Response of the Core when receiving an result from the `SubmitResult` call

**Example:**
```json
{}
```




<a name="SubmitResultRequest"/>

### SubmitResultRequest
Data sent while submitting a a result for a task.
This result can only be called once inside a request from `ListenTask` because of its dependency with the `executionID``

**Example:**
```json
{
  "executionID": "xxxxxx",
  "outputKey": "outputX"
  "outputData": "{\"foo\":\"super result\",\"bar\":true}"
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| executionID | [string](#string) |  | The executionID received from the `listenTask` stream |
| outputKey | [string](#string) |  | The output key defined in your [mesg.yml](/guide/service/service-file.html). |
| outputData | [string](#string) |  | The data of your result encoded in JSON as defined in your [mesg.yml](/guide/service/service-file.html). |




<a name="TaskData"/>

### TaskData
Data sent through the stream from the `ListenTask` API
These data can come as long as the stream stays open. They contains all necessary informations to process a task.
The `executionID` needs to be kept and sent back with the `submitResult` API

**Example:**
```json
{
  "executionID": "xxxxxx",
  "taskKey": "taskX",
  "inputData": "{\"inputX\":\"Hello world!\",\"inputY\":true}"
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| executionID | [string](#string) |  | An unique identifier for the execution you want to submit the result to. |
| taskKey | [string](#string) |  | The task key defined in your [mesg.yml](/guide/service/service-file.html). |
| inputData | [string](#string) |  | The inputs of your tasks encoded in JSON as defined in your [mesg.yml](/guide/service/service-file.html). |





 <!-- end enums -->


