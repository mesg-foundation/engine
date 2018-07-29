



# Data service
<!--
DO NOT EDIT
This file is generated using the ./scripts/build-proto.sh scripts
Please update the github.com/mesg-foundation/core/service/service.proto file
-->



[[toc]]
















## Service


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the service |
| description | [string](#string) |  | Description of the service |
| tasks | [Service.TasksEntry](#Service.TasksEntry) | repeated | Map of tasks that the service can execute |
| events | [Service.EventsEntry](#Service.EventsEntry) | repeated | Map of events that the service can emits |
| dependencies | [Service.DependenciesEntry](#Service.DependenciesEntry) | repeated | Docker dependencies that the service requires |
| configuration | [Dependency](#Dependency) |  | Docker configurations for the service |
| repository | [string](#string) |  | Repository where the source code of this service is accessible |




















## Dependency


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| image | [string](#string) |  | Docker image name or sha |
| volumes | [string](#string) | repeated | List of Docker volumes |
| volumesfrom | [string](#string) | repeated | List of volumes mounted from other dependencies |
| ports | [string](#string) | repeated | List of ports that the container needs to expose |
| command | [string](#string) |  | Command needed to run the container |






## Event


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the event |
| description | [string](#string) |  | Description of the event |
| data | [Event.DataEntry](#Event.DataEntry) | repeated | Map of data associated to this event |






## Event.DataEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Parameter](#Parameter) |  |  |






## Output


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the output |
| description | [string](#string) |  | Description of the output |
| data | [Output.DataEntry](#Output.DataEntry) | repeated | Map of data associated to this output |






## Output.DataEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Parameter](#Parameter) |  |  |






## Parameter


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the parameter |
| description | [string](#string) |  | Description of the parameter |
| type | [string](#string) |  | Type of the parameter `String`, `Number`, `Boolean` or `Object` |
| optional | [bool](#bool) |  | Mark this parameter as optional |








## Service.DependenciesEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Dependency](#Dependency) |  |  |






## Service.EventsEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Event](#Event) |  |  |






## Service.TasksEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Task](#Task) |  |  |






## Task


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the task |
| description | [string](#string) |  | Description of the task |
| inputs | [Task.InputsEntry](#Task.InputsEntry) | repeated | Map of inputs that can be given for this task |
| outputs | [Task.OutputsEntry](#Task.OutputsEntry) | repeated | Map of outputs that the task can returns as result |






## Task.InputsEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Parameter](#Parameter) |  |  |






## Task.OutputsEntry


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Output](#Output) |  |  |





