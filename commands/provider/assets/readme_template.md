# {{.Name}}

{{.Description}}
{{if .Repository}}
```bash
mesg-core service deploy {{.Repository}}
```
{{end}}
{{if .Events}}# Events
{{range $key, $event := .Events}}
## {{$key}}

Event key: `{{$key}}`

{{$event.Description}}

{{if $event.Data}}| **Key** | **Type** | **Description** |
| --- | --- | --- |
{{range $dataKey, $data := $event.Data}}| **{{$dataKey}}** | `{{$data.Type}}` | {{$data.Description}} |
{{end}}{{end}}{{end}}{{end}}
{{if .Tasks}}
# Tasks
{{range $key, $task := .Tasks}}
## {{$key}}

Task key: `{{$key}}`

{{$task.Description}}

{{if $task.Inputs}}### Inputs

| **Key** | **Type** | **Description** |
| --- | --- | --- |
{{range $inputKey, $input := $task.Inputs}}| **{{$inputKey}}** | `{{$input.Type}}` | {{$input.Description}} |
{{end}}{{end}}

{{if $task.Outputs}}### Outputs

{{range $outputKey, $output := $task.Outputs}}##### {{$outputKey}}

Output key: `{{$outputKey}}`

{{$output.Description}}

| **Key** | **Type** | **Description** |
| --- | --- | --- |
{{range $outputKey, $output := $output.Data}}| **{{$outputKey}}** | `{{$output.Type}}` | {{$output.Description}} |
{{end}}
{{end}}{{end}}

{{end}}{{end}}
