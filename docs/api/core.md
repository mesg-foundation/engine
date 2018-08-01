




# Core API
<!--
DO NOT EDIT
This file is generated using the ./scripts/build-proto.sh scripts
Please update the Core file
-->

Main API to interact with MESG Core functionalities. It can be consumed by any applications or any tools that wish to interact with MESG Core.
It is actually used by the MESG CLI and MESG Application libraries.

Services should not use this API but use the [Service API](/api/service.html).

[[toc]]


## ListenEvent

Subscribe to a stream that listens for events from a service.

<tabs>
<tab title="Request">

























#### ListenEventRequest
Request's data of the `ListenEvent` stream API.

**Example**
```json
{
  "serviceID":   "__SERVICE_ID__",
  "eventFilter": "__EVENT_KEY_TO_MATCH__"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| serviceID | [string](#string) | Service ID. Generated when using the `DeployService` API. |
| eventFilter | [string](#string) | __Optional.__ Event's key to filter. The event have to match this key. The default is `*` and matches any event. |

















</tab>

<tab title="Reply">











#### EventData
Data receive from the stream of the `ListenEvent` API.
Will be received over time as long as the stream is opened.

**Example**
```json
{
  "eventKey":  "__EVENT_KEY__",
  "eventData": "{\"foo\":\"bar\"}"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| eventKey | [string](#string) | Event's key. |
| eventData | [string](#string) | Event's data encoded in JSON. |































</tab>
</tabs>

## ListenResult

Subscribe to the stream that listens for task's results of a service.

<tabs>
<tab title="Request">



























#### ListenResultRequest
Request's data of the `ListenResult` stream API.

**Example**
```json
{
  "serviceID":     "__SERVICE_ID__",
  "taskFilter":    "__TASK_KEY_TO_MATCH__",
  "outputFilter":  "__OUTPUT_KEY_TO_MATCH__"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| serviceID | [string](#string) | Service ID. Generated when using the `DeployService` API. |
| taskFilter | [string](#string) | __Optional.__  Task's key to filter. The task have to match this key. The default is `*` and matches any task. |
| outputFilter | [string](#string) | __Optional.__ Output's key of the task to filter. The task have to return this output's key. The default is `*` and matches any output. |















</tab>

<tab title="Reply">





























#### ResultData
Data receive from the stream of the `ListenResult` API.
Will be received over time as long as the stream is opened.

**Example**
```json
{
  "executionID": "__EXECUTION_ID__",
  "taskKey":     "__TASK_KEY__",
  "outputKey":   "__OUTPUT_KEY__",
  "outputData":  "{\"foo\":\"bar\"}"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| executionID | [string](#string) | Unique identifier of the execution. |
| taskKey | [string](#string) | Key of the executed task. |
| outputKey | [string](#string) | Output's key the task returned. |
| outputData | [string](#string) | Output's data the task returned encoded in JSON. |













</tab>
</tabs>

## ExecuteTask

Execute a task of a service through the [Core](/guide/start-here/core.html).

<tabs>
<tab title="Request">















#### ExecuteTaskRequest
Request's data of the `ExecuteTask` API.

**Example**
```json
{
  "serviceID": "__SERVICE_ID__",
  "taskKey":   "__TASK_KEY__",
  "inputData": "{\"foo\":\"bar\"}"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| serviceID | [string](#string) | Service ID. Generated when using the `DeployService` API. |
| taskKey | [string](#string) | Task's key to execute. |
| inputData | [string](#string) | Inputs of the task to execute encoded in JSON. |



























</tab>

<tab title="Reply">













#### ExecuteTaskReply
Reply's data of the `ExecuteTask` API.

**Example**
```json
{
  "executionID": "__EXECUTION_ID__"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| executionID | [string](#string) | Unique identifier of the execution. |





























</tab>
</tabs>

## StartService

Start a service. The service have to be already deployed on the [Core](/guide/start-here/core.html).

<tabs>
<tab title="Request">

































#### StartServiceRequest
Request's data of the `StartService` API.

**Example**
```json
{
  "serviceID": "__SERVICE_ID__"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| serviceID | [string](#string) | Service ID. Generated when using the `DeployService` API. |









</tab>

<tab title="Reply">































#### StartServiceReply
Reply's data of the `StartService` API.











</tab>
</tabs>

## StopService

Stop a service. The service have to be already deployed on the [Core](/guide/start-here/core.html).

<tabs>
<tab title="Request">





































#### StopServiceRequest
Request's data of the `StopService` API.

**Example**
```json
{
  "serviceID": "__SERVICE_ID__"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| serviceID | [string](#string) | Service ID. Generated when using the `DeployService` API. |





</tab>

<tab title="Reply">



































#### StopServiceReply
Reply's data of the `StopService` API.







</tab>
</tabs>

## DeployService

Deploy a service to the [Core](/guide/start-here/core.html). This will give you an unique identifier to use to interact with the service.

<tabs>
<tab title="Request">









#### DeployServiceRequest
Request's data of the `DeployService` API.

**Example**
```json
{
  "service": {
    "name": "serviceX",
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
| service | [service.Service](#service.Service) | Service's definition to deploy. [Details here](/api/service-type.html) |

































</tab>

<tab title="Reply">







#### DeployServiceReply
Reply's data of the `DeployService` API.

**Example**
```json
{
  "serviceID": "__SERVICE_ID__"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| serviceID | [string](#string) | The generated identifier of the deployed service. Use this ID with other APIs. |



































</tab>
</tabs>

## DeleteService

Delete a service from Core. This function only delete a deployed service in the [Core](/guide/start-here/core.html). If the service's code is on your computer, it will not delete its source code.

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
| serviceID | [string](#string) | Service ID. Generated when using the `DeployService` API. |





































</tab>

<tab title="Reply">



#### DeleteServiceReply
Reply's data  of the `DeleteService` API.







































</tab>
</tabs>

## ListServices

List all services already deployed in the [Core](/guide/start-here/core.html).

<tabs>
<tab title="Request">























#### ListServicesRequest
Request's data of the `ListServices` API.



















</tab>

<tab title="Reply">





















#### ListServicesReply
Reply's data of the `ListServices` API.

**Example**
```json
[{
  "service": {
    "name": "serviceX",
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
| services | [service.Service](#service.Service)[] | List of services' definition previously deployed. [Details here](/api/service-type.html) |





















</tab>
</tabs>

## GetService

Get the definition of an already deployed service from its ID.

<tabs>
<tab title="Request">



















#### GetServiceRequest
Request's data of the `GetService` API.

**Example**
```json
{
  "serviceID": "__SERVICE_ID__"
}
```


| Field | Type | Description |
| ----- | ---- | ----------- |
| serviceID | [string](#string) | Service ID. Generated when using the `DeployService` API. |























</tab>

<tab title="Reply">

















#### GetServiceReply
Reply's data of the `GetService` API.

**Example**
```json
{
  "service": {
    "name": "serviceX",
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
| service | [service.Service](#service.Service) | Service's definition. [Details here](/api/service-type.html) |

























</tab>
</tabs>



