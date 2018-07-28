

# Core API

## APIs



| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| ListenEvent | [ListenEventRequest](#ListenEventRequest) | [EventData](#ListenEventRequest) |  |
| ExecuteTask | [ExecuteTaskRequest](#ExecuteTaskRequest) | [ExecuteTaskReply](#ExecuteTaskRequest) |  |
| ListenResult | [ListenResultRequest](#ListenResultRequest) | [ResultData](#ListenResultRequest) |  |
| StartService | [StartServiceRequest](#StartServiceRequest) | [StartServiceReply](#StartServiceRequest) |  |
| StopService | [StopServiceRequest](#StopServiceRequest) | [StopServiceReply](#StopServiceRequest) |  |
| DeployService | [DeployServiceRequest](#DeployServiceRequest) | [DeployServiceReply](#DeployServiceRequest) |  |
| DeleteService | [DeleteServiceRequest](#DeleteServiceRequest) | [DeleteServiceReply](#DeleteServiceRequest) |  |
| ListServices | [ListServicesRequest](#ListServicesRequest) | [ListServicesReply](#ListServicesRequest) |  |
| GetService | [GetServiceRequest](#GetServiceRequest) | [GetServiceReply](#GetServiceRequest) |  |



## Types




<a name="DeleteServiceReply"/>

### DeleteServiceReply







<a name="DeleteServiceRequest"/>

### DeleteServiceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceID | [string](#string) |  |  |






<a name="DeployServiceReply"/>

### DeployServiceReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceID | [string](#string) |  |  |






<a name="DeployServiceRequest"/>

### DeployServiceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service | [service.Service](#service.Service) |  |  |






<a name="EventData"/>

### EventData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| eventKey | [string](#string) |  |  |
| eventData | [string](#string) |  |  |






<a name="ExecuteTaskReply"/>

### ExecuteTaskReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| executionID | [string](#string) |  |  |






<a name="ExecuteTaskRequest"/>

### ExecuteTaskRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceID | [string](#string) |  |  |
| taskKey | [string](#string) |  |  |
| inputData | [string](#string) |  |  |






<a name="GetServiceReply"/>

### GetServiceReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service | [service.Service](#service.Service) |  |  |






<a name="GetServiceRequest"/>

### GetServiceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceID | [string](#string) |  |  |






<a name="ListServicesReply"/>

### ListServicesReply



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| services | [service.Service](#service.Service) | repeated |  |






<a name="ListServicesRequest"/>

### ListServicesRequest







<a name="ListenEventRequest"/>

### ListenEventRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceID | [string](#string) |  |  |
| eventFilter | [string](#string) |  |  |






<a name="ListenResultRequest"/>

### ListenResultRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceID | [string](#string) |  |  |
| taskFilter | [string](#string) |  |  |
| outputFilter | [string](#string) |  |  |






<a name="ResultData"/>

### ResultData



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| executionID | [string](#string) |  |  |
| taskKey | [string](#string) |  |  |
| outputKey | [string](#string) |  |  |
| outputData | [string](#string) |  |  |






<a name="StartServiceReply"/>

### StartServiceReply







<a name="StartServiceRequest"/>

### StartServiceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceID | [string](#string) |  |  |






<a name="StopServiceReply"/>

### StopServiceReply







<a name="StopServiceRequest"/>

### StopServiceRequest



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceID | [string](#string) |  |  |







 <!-- end enums -->

 <!-- end HasExtensions -->


