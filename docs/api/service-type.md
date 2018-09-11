



# Data core
<!--
DO NOT EDIT
This file is generated using the ./scripts/build-proto.sh scripts
Please update the github.com/mesg-foundation/core/interface/grpc/core/service.proto file
-->



[[toc]]
















#### Service
This is the definition of a MESG Service.


| Field | Type | Description |
| ----- | ---- | ----------- |
| ID | [string](#string) | Service's unique id service hash. |
| name | [string](#string) | Service's name. |
| description | [string](#string) | Service's description. |
| tasks | [Service.TasksEntry](#core.Service.TasksEntry)[] | The list of tasks this service can execute. |
| events | [Service.EventsEntry](#core.Service.EventsEntry)[] | The list of events this service can emit. |
| dependencies | [Service.DependenciesEntry](#core.Service.DependenciesEntry)[] | The Docker dependencies this service requires. |
| configuration | [Dependency](#core.Dependency) | Service's Docker configuration. |
| repository | [string](#string) | Service's repository that contain its source code. |





















#### Dependency
A dependency is a configuration of an other Docker container that runs separately from the service.


| Field | Type | Description |
| ----- | ---- | ----------- |
| image | [string](#string) | Image's name of the Docker. |
| volumes | [string](#string)[] | List of volumes. |
| volumesfrom | [string](#string)[] | List of volumes mounted from other dependencies. |
| ports | [string](#string)[] | List of ports the container exposes. |
| command | [string](#string) | Command to run the container. |







#### Event
Events are emitted by the service whenever the service wants.
TODO(ilgooz) remove key, serviceName fields when Event type crafted manually.


| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) | Event's key. |
| name | [string](#string) | Event's name. |
| description | [string](#string) | Event's description. |
| serviceName | [string](#string) | Event's service name. |
| data | [Event.DataEntry](#core.Event.DataEntry)[] | List of data of this event. |







#### Event.DataEntry



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) |  |
| value | [Parameter](#core.Parameter) |  |







#### Output
A output is the data a task must return.
TODO(ilgooz) remove key, taskKey, serviceName fields when Output type crafted manually.


| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) | Output's key. |
| name | [string](#string) | Output's name. |
| description | [string](#string) | Output's description. |
| taskKey | [string](#string) | Output's task key. |
| serviceName | [string](#string) | Output's service name. |
| data | [Output.DataEntry](#core.Output.DataEntry)[] | List of data of this output. |







#### Output.DataEntry



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) |  |
| value | [Parameter](#core.Parameter) |  |







#### Parameter
A parameter is the definition of a specific value.


| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [string](#string) | Parameter's name. |
| description | [string](#string) | Parameter's description. |
| type | [string](#string) | Parameter's type: `String`, `Number`, `Boolean` or `Object`. |
| optional | [bool](#bool) | Set the parameter as optional. |









#### Service.DependenciesEntry



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) |  |
| value | [Dependency](#core.Dependency) |  |







#### Service.EventsEntry



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) |  |
| value | [Event](#core.Event) |  |







#### Service.TasksEntry



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) |  |
| value | [Task](#core.Task) |  |







#### Task
A task is a function that requires inputs and returns output.
TODO(ilgooz) remove key, serviceName fields when Task type crafted manually.


| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) | Task's key. |
| name | [string](#string) | Task's name. |
| description | [string](#string) | Task's description. |
| serviceName | [string](#string) | Task's service name. |
| inputs | [Task.InputsEntry](#core.Task.InputsEntry)[] | List inputs of this task. |
| outputs | [Task.OutputsEntry](#core.Task.OutputsEntry)[] | List of outputs this task can return. |







#### Task.InputsEntry



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) |  |
| value | [Parameter](#core.Parameter) |  |







#### Task.OutputsEntry



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) |  |
| value | [Output](#core.Output) |  |






