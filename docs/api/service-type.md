




# Service Data
<!--
DO NOT EDIT
This file is generated using the ./scripts/build-proto.sh scripts
Please update the github.com/mesg-foundation/core/service/service.proto file
-->



<a name="Dependency"/>

### Dependency



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| image | [string](#string) |  | Docker image name or sha |
| volumes | [string](#string) | repeated | List of Docker volumes |
| volumesfrom | [string](#string) | repeated | List of volumes mounted from other dependencies |
| ports | [string](#string) | repeated | List of ports that the container needs to expose |
| command | [string](#string) |  | Command needed to run the container |




<a name="Event"/>

### Event



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the event |
| description | [string](#string) |  | Description of the event |
| data | [Event.DataEntry](#service.Event.DataEntry) | repeated | Map of data associated to this event |




<a name="DataEntry"/>

### Event.DataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Parameter](#service.Parameter) |  |  |




<a name="Output"/>

### Output



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the output |
| description | [string](#string) |  | Description of the output |
| data | [Output.DataEntry](#service.Output.DataEntry) | repeated | Map of data associated to this output |




<a name="DataEntry"/>

### Output.DataEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Parameter](#service.Parameter) |  |  |




<a name="Parameter"/>

### Parameter



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the parameter |
| description | [string](#string) |  | Description of the parameter |
| type | [string](#string) |  | Type of the parameter `String`, `Number`, `Boolean` or `Object` |
| optional | [bool](#bool) |  | Mark this parameter as optional |




<a name="Service"/>

### Service



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the service |
| description | [string](#string) |  | Description of the service |
| tasks | [Service.TasksEntry](#service.Service.TasksEntry) | repeated | Map of tasks that the service can execute |
| events | [Service.EventsEntry](#service.Service.EventsEntry) | repeated | Map of events that the service can emits |
| dependencies | [Service.DependenciesEntry](#service.Service.DependenciesEntry) | repeated | Docker dependencies that the service requires |
| configuration | [Dependency](#service.Dependency) |  | Docker configurations for the service |
| repository | [string](#string) |  | Repository where the source code of this service is accessible |




<a name="DependenciesEntry"/>

### Service.DependenciesEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Dependency](#service.Dependency) |  |  |




<a name="EventsEntry"/>

### Service.EventsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Event](#service.Event) |  |  |




<a name="TasksEntry"/>

### Service.TasksEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Task](#service.Task) |  |  |




<a name="Task"/>

### Task



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name of the task |
| description | [string](#string) |  | Description of the task |
| inputs | [Task.InputsEntry](#service.Task.InputsEntry) | repeated | Map of inputs that can be given for this task |
| outputs | [Task.OutputsEntry](#service.Task.OutputsEntry) | repeated | Map of outputs that the task can returns as result |




<a name="InputsEntry"/>

### Task.InputsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Parameter](#service.Parameter) |  |  |




<a name="OutputsEntry"/>

### Task.OutputsEntry



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| key | [string](#string) |  |  |
| value | [Output](#service.Output) |  |  |





 <!-- end enums -->


