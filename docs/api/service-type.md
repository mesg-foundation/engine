



# Data service
<!--
DO NOT EDIT
This file is generated using the ./scripts/build-proto.sh scripts
Please update the github.com/mesg-foundation/core/service/service.proto file
-->



[[toc]]
















#### Service



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [string](#string) | Name of the service |
| description | [string](#string) | Description of the service |
| tasks | [Service.TasksEntry](#service.Service.TasksEntry)[] | Map of tasks that the service can execute |
| events | [Service.EventsEntry](#service.Service.EventsEntry)[] | Map of events that the service can emits |
| dependencies | [Service.DependenciesEntry](#service.Service.DependenciesEntry)[] | Docker dependencies that the service requires |
| configuration | [Dependency](#service.Dependency) | Docker configurations for the service |
| repository | [string](#string) | Repository where the source code of this service is accessible |





















#### Dependency



| Field | Type | Description |
| ----- | ---- | ----------- |
| image | [string](#string) | Docker image name or sha |
| volumes | [string](#string)[] | List of Docker volumes |
| volumesfrom | [string](#string)[] | List of volumes mounted from other dependencies |
| ports | [string](#string)[] | List of ports that the container needs to expose |
| command | [string](#string) | Command needed to run the container |







#### Event



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [string](#string) | Name of the event |
| description | [string](#string) | Description of the event |
| data | [Event.DataEntry](#service.Event.DataEntry)[] | Map of data associated to this event |







#### Event.DataEntry



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) |  |
| value | [Parameter](#service.Parameter) |  |







#### Output



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [string](#string) | Name of the output |
| description | [string](#string) | Description of the output |
| data | [Output.DataEntry](#service.Output.DataEntry)[] | Map of data associated to this output |







#### Output.DataEntry



| Field | Type | Description |
| ----- | ---- | ----------- |
| key | [string](#string) |  |
| value | [Parameter](#service.Parameter) |  |







#### Parameter



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [string](#string) | Name of the parameter |
| description | [string](#string) | Description of the parameter |
| type | [string](#string) | Type of the parameter `String`, `Number`, `Boolean` or `Object` |
| optional | [bool](#bool) | Mark this parameter as optional |









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



| Field | Type | Description |
| ----- | ---- | ----------- |
| name | [string](#string) | Name of the task |
| description | [string](#string) | Description of the task |
| inputs | [Task.InputsEntry](#service.Task.InputsEntry)[] | Map of inputs that can be given for this task |
| outputs | [Task.OutputsEntry](#service.Task.OutputsEntry)[] | Map of outputs that the task can returns as result |







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






