


## Types




<a name="Dependency"/>

### Dependency



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| image | [string](#string) |  |  |
| volumes | [string](#string) | repeated |  |
| volumesfrom | [string](#string) | repeated |  |
| ports | [string](#string) | repeated |  |
| command | [string](#string) |  |  |






<a name="Event"/>

### Event



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| data | [Event.DataEntry](#service.Event.DataEntry) | repeated |  |






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
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| data | [Output.DataEntry](#service.Output.DataEntry) | repeated |  |






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
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| type | [string](#string) |  |  |
| optional | [bool](#bool) |  |  |






<a name="Service"/>

### Service



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| tasks | [Service.TasksEntry](#service.Service.TasksEntry) | repeated |  |
| events | [Service.EventsEntry](#service.Service.EventsEntry) | repeated |  |
| dependencies | [Service.DependenciesEntry](#service.Service.DependenciesEntry) | repeated |  |
| configuration | [Dependency](#service.Dependency) |  |  |
| repository | [string](#string) |  |  |






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
| name | [string](#string) |  |  |
| description | [string](#string) |  |  |
| inputs | [Task.InputsEntry](#service.Task.InputsEntry) | repeated |  |
| outputs | [Task.OutputsEntry](#service.Task.OutputsEntry) | repeated |  |






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

 <!-- end HasExtensions -->


