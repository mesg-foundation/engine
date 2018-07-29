




# Core API
<!--
DO NOT EDIT
This file is generated using the ./scripts/build-proto.sh scripts
Please update the Core file
-->

API accessible for anyone, it can be consumed either by an application or any tool that wishes to connect to MESG.
It is actually used by the MESG CLI.

Services should not try to access this API

[[toc]]


## ListenEvent

Subscribe to the stream that will receive events from a service

<tabs>
<tab title="Request">
























### Request

#### ListenEventRequest
Data sent to connect to the `ListenEvent` stream API

**Example**
```json
{
  "serviceID": "xxxx",
  "eventFilter": "*"
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceID | [string](#string) |  | The ID that references your service. Generated while using the `DeployService` API. |
| eventFilter | [string](#string) |  | The key of the event you want to filter from the service. The default `"*"` will listen everything. |

















</tab>

<tab title="Response">










### Response

#### EventData
Data sent through the stream from the `ListenEvent` API
These data can come as long as the stream stays open.

**Example**
```json
{
  "eventKey": "xxxx",
  "eventData": "{\"foo\":\"bar\"}"
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| eventKey | [string](#string) |  |  |
| eventData | [string](#string) |  |  |































</tab>
</tabs>

## ExecuteTask

Let you to execute a task of a service through the [Core](/guide/start-here/core.html)

<tabs>
<tab title="Request">














### Request

#### ExecuteTaskRequest
Payload sent when you want to execute a task of a service

**Example**
```json
{
  "serviceID": "xxxx",
  "taskKey": "myTaskX",
  "inputData": "{\"foo\":\"bar\"}"
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceID | [string](#string) |  | The ID that references your service. Generated while using the `DeployService` API. |
| taskKey | [string](#string) |  | The key of the task you want to execute from the service. |
| inputData | [string](#string) |  | The inputs for the tasks you want to execute encoded in JSON. |



























</tab>

<tab title="Response">












### Response

#### ExecuteTaskReply
Reply of the [Core](/guide/start-here/core.html) when calling the `ExecuteTask` API

**Example**
```json
{
  "executionID": "xxx"
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| executionID | [string](#string) |  | The unique identifier for this execution that let you track the result. |





























</tab>
</tabs>

## ListenResult

Subscribe to the stream that will receive results of a task of a service

<tabs>
<tab title="Request">


























### Request

#### ListenResultRequest
Data sent to connect to the `ListenResult` stream API

**Example**
```json
{
  "serviceID": "xxxx",
  "taskFilter": "*",
  "outputFilter": "*"
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceID | [string](#string) |  | The ID that references your service. Generated while using the `DeployService` API. |
| taskFilter | [string](#string) |  | The key of the task you want to filter from the service. The default `"*"` will listen everything. |
| outputFilter | [string](#string) |  | The key of the output you want to filter from the service. The default `"*"` will listen everything. |















</tab>

<tab title="Response">




























### Response

#### ResultData
Data sent to the `ListenResult` stream that contains all informations of a result execution

**Example**
```json
{
  "executionID": "xxx",
  "taskKey": "taskX",
  "outputKey": "outputX",
  "outputData": "{\"foo\":\"bar\"}"
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| executionID | [string](#string) |  | The unique identifier of your execution |
| taskKey | [string](#string) |  | The key of the task executed |
| outputKey | [string](#string) |  | The key of the output the task returned |
| outputData | [string](#string) |  | The data of the output the task returned encoded in JSON |













</tab>
</tabs>

## StartService

Start a service. This service needs to be deployed already in the [Core](/guide/start-here/core.html)

<tabs>
<tab title="Request">
































### Request

#### StartServiceRequest
Payload necessary to start a service

**Example**
```json
{
  "serviceID": "xxxx"
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceID | [string](#string) |  | The ID that references your service. Generated while using the `DeployService` API. |









</tab>

<tab title="Response">






























### Response

#### StartServiceReply
Reply of the [Core](/guide/start-here/core.html) whan starting a Service

**Example**
```json
{}
```











</tab>
</tabs>

## StopService

Stop a service. This service needs to be deployed already in the [Core](/guide/start-here/core.html)

<tabs>
<tab title="Request">




































### Request

#### StopServiceRequest
Payload necessary to stop a service

**Example**
```json
{
  "serviceID": "xxxx"
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceID | [string](#string) |  | The ID that references your service. Generated while using the `DeployService` API. |





</tab>

<tab title="Response">


































### Response

#### StopServiceReply
Reply of the [Core](/guide/start-here/core.html) whan stopping a Service

**Example**
```json
{}
```







</tab>
</tabs>

## DeployService

Deploy a new service to the [Core](/guide/start-here/core.html). This will give you an unique identifier to use your service

<tabs>
<tab title="Request">








### Request

#### DeployServiceRequest
Data sent while deploying a new Service to the [Core](/guide/start-here/core.html)

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


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service | [service.Service](#service.Service) |  | Data of the service you want to deploy. [details here](/api/service-type.html) |

































</tab>

<tab title="Response">






### Response

#### DeployServiceReply
Reply of the [Core](/guide/start-here/core.html) whan deploying a new Service

**Example**
```json
{
  "serviceID": "xxx"
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceID | [string](#string) |  | The generated ID for your service. You can reuse this ID to access many others APIs |



































</tab>
</tabs>

## DeleteService

Delete a service. This function will only delete the service deployed in the [Core](/guide/start-here/core.html). If the service code is on your computer, this call will not delete your source code

<tabs>
<tab title="Request">




### Request

#### DeleteServiceRequest
Payload necessary to delete a service

**Example**
```json
{
  "serviceID": "xxxx"
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceID | [string](#string) |  | The ID that references your service. Generated while using the `DeployService` API. |





































</tab>

<tab title="Response">


### Response

#### DeleteServiceReply
Reply of the [Core](/guide/start-here/core.html) whan deleting a Service

**Example**
```json
{}
```







































</tab>
</tabs>

## ListServices

List all the services already deployed in the [Core](/guide/start-here/core.html)

<tabs>
<tab title="Request">






















### Request

#### ListServicesRequest
Payload necessary to list all the deployed services

**Example**
```json
{}
```



















</tab>

<tab title="Response">




















### Response

#### ListServicesReply
Result from the [Core](/guide/start-here/core.html) when calling the list of services deployed

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


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| services | [service.Service](#service.Service) | repeated | List of data of the deployed services. [details here](/api/service-type.html) |





















</tab>
</tabs>

## GetService

Get an already deployed service based on its ID

<tabs>
<tab title="Request">


















### Request

#### GetServiceRequest
Payload necessary to get the details of deployed service

**Example**
```json
{
  "serviceID": "xxxx"
}
```


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| serviceID | [string](#string) |  | The ID that references your service. Generated while using the `DeployService` API. |























</tab>

<tab title="Response">
















### Response

#### GetServiceReply
Result from the [Core](/guide/start-here/core.html) when calling the `GetService` API

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


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| service | [service.Service](#service.Service) |  | Data of the service you requested. [details here](/api/service-type.html) |

























</tab>
</tabs>



