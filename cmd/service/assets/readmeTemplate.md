# {{.Name}}

{{.Description}}

```bash
mesg-core service deploy __REPLACE_BY_YOUR_REPOSITORY__
```

{{if .Events}}
## Events
{{range $key, $event := .Events}}
### {{$key}}

Event key: **{{$key}}**

{{$event.Description}}

{{if $event.Data}}| **key** | **type** | **description** |
| --- | --- | --- |
{{range $dataKey, $data := $event.Data}}| **{{$dataKey}}** | `{{$data.Type}}` | {{$data.Description}} |
{{end}}{{end}}{{end}}{{end}}
{{if .Tasks}}
## Tasks
{{range $key, $task := .Tasks}}
### {{$key}}

Task key: **{{$key}}**

{{$task.Description}}

{{if $task.Inputs}}#### Inputs

| **key** | **type** | **description** |
| --- | --- | --- |
{{range $inputKey, $input := $task.Inputs}}| **{{$inputKey}}** | `{{$input.Type}}` | {{$input.Description}} |
{{end}}{{end}}

{{if $task.Outputs}}#### Outputs

{{range $outputKey, $output := $task.Outputs}}##### {{$outputKey}}

Output key: **{{$outputKey}}**

{{$output.Description}}

| **key** | **type** | **description** |
| --- | --- | --- |
{{range $outputKey, $output := $output.Data}}| **{{$outputKey}}** | `{{$output.Type}}` | {{$output.Description}} |
{{end}}
{{end}}{{end}}

{{end}}{{end}}
