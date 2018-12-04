# {{.Name}}

{{.Description}}
{{if .Repository}}
# Contents

- [Installation](#Installation)
- [Definitions](#Definitions)
  {{if .Events}}- [Events](#Events){{range $key, $event := .Events}}
    - [{{or $event.Name $key}}](#{{or $event.Name $key | anchorEncode}}){{end}}{{end}}
  {{if .Tasks}}- [Tasks](#Tasks){{range $key, $task := .Tasks}}
    - [{{or $task.Name $key}}](#{{or $task.Name $key | anchorEncode}}){{end}}{{end}}

# Installation

## MESG Core

This service requires [MESG Core](https://github.com/mesg-foundation/core) to be install.

You can install MESG Core by running the following command or [follow the install guide](https://docs.mesg.tech/guide/start-here/installation.html).

```bash
bash <(curl -fsSL https://mesg.com/install)
```
```bash
mesg-core service deploy {{.Repository}}
```
{{end}}
# Definitions

{{if .Events}}# Events
{{range $key, $event := .Events}}
## {{if eq $event.Name ""}}{{$key}}{{ else }}{{ $event.Name }}{{end}}

Event key: `{{$key}}`

{{$event.Description}}

{{if $event.Data}}| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
{{range $dataKey, $data := $event.Data}}| **{{if eq $data.Name ""}}{{$dataKey}}{{ else }}{{ $data.Name }}{{end}}** | `{{$dataKey}}` | `{{$data.Type}}` | {{if $data.Optional}}**`optional`** {{end}}{{$data.Description}} |
{{end}}{{end}}{{end}}{{end}}
{{if .Tasks}}# Tasks
{{range $key, $task := .Tasks}}
## {{if eq $task.Name ""}}{{$key}}{{ else }}{{ $task.Name }}{{end}}

Task key: `{{$key}}`

{{$task.Description}}

{{if $task.Inputs}}### Inputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
{{range $inputKey, $input := $task.Inputs}}| **{{if eq $input.Name ""}}{{$inputKey}}{{ else }}{{ $input.Name }}{{end}}** | `{{$inputKey}}` | `{{$input.Type}}` | {{if $input.Optional}}**`optional`** {{end}}{{$input.Description}} |
{{end}}{{end}}
{{if $task.Outputs}}### Outputs

{{range $outputKey, $output := $task.Outputs}}#### {{if eq $output.Name ""}}{{$outputKey}}{{ else }}{{ $output.Name }}{{end}}

Output key: `{{$outputKey}}`

{{$output.Description}}

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
{{range $outputKey, $output := $output.Data}}| **{{if eq $output.Name ""}}{{$outputKey}}{{ else }}{{ $output.Name }}{{end}}** | `{{$outputKey}}` | `{{$output.Type}}` | {{if $output.Optional}}**`optional`** {{end}}{{$output.Description}} |
{{end}}
{{end}}{{end}}{{end}}{{end}}
