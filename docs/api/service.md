

# Service API

## APIs



| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| EmitEvent | [EmitEventRequest](#EmitEventRequest) | [EmitEventReply](#EmitEventRequest) |  |
| ListenTask | [ListenTaskRequest](#ListenTaskRequest) | [TaskData](#ListenTaskRequest) |  |
| SubmitResult | [SubmitResultRequest](#SubmitResultRequest) | [SubmitResultReply](#SubmitResultRequest) |  |



## Types




<a name="EmitEventReply"/>

### EmitEventReply







<a name="EmitEventRequest"/>

### EmitEventRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token | [string](#string) |  |  |
| eventKey | [string](#string) |  |  |
| eventData | [string](#string) |  |  |






<a name="ListenTaskRequest"/>

### ListenTaskRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| token | [string](#string) |  |  |






<a name="SubmitResultReply"/>

### SubmitResultReply







<a name="SubmitResultRequest"/>

### SubmitResultRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| executionID | [string](#string) |  |  |
| outputKey | [string](#string) |  |  |
| outputData | [string](#string) |  |  |






<a name="TaskData"/>

### TaskData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| executionID | [string](#string) |  |  |
| taskKey | [string](#string) |  |  |
| inputData | [string](#string) |  |  |







 <!-- end enums -->

 <!-- end HasExtensions -->


