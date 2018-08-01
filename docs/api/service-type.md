



# Data service
<!--
DO NOT EDIT
This file is generated using the ./scripts/build-proto.sh scripts
Please update the github.com/mesg-foundation/core/service/service.proto file
-->



[[toc]]
















#### Service
This is the definition of a MESG Service.


| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [string](#string) | The service's name. |
| description | [string](#string) | The service's description. |
| tasks | [Service.TasksEntry](#service.Service.TasksEntry)[] | The list of tasks this service can execute. |
| events | [Service.EventsEntry](#service.Service.EventsEntry)[] | The list of events this service can emit. |
| dependencies | [Service.DependenciesEntry](#service.Service.DependenciesEntry)[] | The Docker dependencies this service requires. |
| configuration | [Dependency](#service.Dependency) | The service's Docker configuration. |
| repository | [string](#string) | The service's repository that contain its source code. |





















#### Dependency
A dependency is a configuration of an other Docker container that runs separately from the service.


| Field | Type | Description |
| ----- | ---- | ----------- |
| image | [string](#string) | Image's name of the Docker |
| volumes | [string](#string)[] | List of volumes |
| volumesfrom | [string](#string)[] | List of volumes mounted from other dependencies |
| ports | [string](#string)[] | List of ports the container exposes |
| command | [string](#string) | Command to run the container |







#### Event
Events are emitted by the service whenever the service wants.


| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [string](#string) | The event's name. |
| description | [string](#string) | The event's description. |
| data | [Event.DataEntry](#service.Event.DataEntry)[] | The list of data of this event. |







#### Event.DataEntry



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) |  |
| value | [Parameter](#service.Parameter) |  |







#### Output
A output is the data a task must return.


| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [string](#string) | The output's name. |
| description | [string](#string) | the output's description. |
| data | [Output.DataEntry](#service.Output.DataEntry)[] | The list of data of this output. |







#### Output.DataEntry



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) |  |
| value | [Parameter](#service.Parameter) |  |







#### Parameter
A parameter is the definition of a specific value.


| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [string](#string) | The parameter's name. |
| description | [string](#string) | The parameter's description. |
| type | [string](#string) | The parameter's type: `String`, `Number`, `Boolean` or `Object`. |
| optional | [bool](#bool) | Set the parameter as optional. |









#### Service.DependenciesEntry



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) |  |
| value | [Dependency](#service.Dependency) |  |







#### Service.EventsEntry



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) |  |
| value | [Event](#service.Event) |  |







#### Service.TasksEntry



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) |  |
| value | [Task](#service.Task) |  |







#### Task
A task is a function that requires inputs and returns output.


| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [string](#string) | The task's name. |
| description | [string](#string) | The task's description. |
| inputs | [Task.InputsEntry](#service.Task.InputsEntry)[] | The list inputs of this task. |
| outputs | [Task.OutputsEntry](#service.Task.OutputsEntry)[] | The list of outputs this task can return. |







#### Task.InputsEntry



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) |  |
| value | [Parameter](#service.Parameter) |  |







#### Task.OutputsEntry



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) |  |
| value | [Output](#service.Output) |  |






