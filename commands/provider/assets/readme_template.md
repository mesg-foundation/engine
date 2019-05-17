# {{.Name}} {{if .Sid}}(ID: {{.Sid}}){{end}}

{{.Description}}

# Contents

- [Installation](#Installation)
- [Definitions](#Definitions)
  {{if .Events}}- [Events](#Events){{range $key, $event := .Events}}
    - [{{or $event.Name $key}}](#{{or $event.Name $key | anchorEncode}}){{end}}{{end}}
  {{if .Tasks}}- [Tasks](#Tasks){{range $key, $task := .Tasks}}
    - [{{or $task.Name $key}}](#{{or $task.Name $key | anchorEncode}}){{end}}{{end}}

# Installation

## MESG Core

This service requires [MESG Core](https://github.com/mesg-foundation/core) to be installed first.

You can install MESG Core by running the following command or [follow the installation guide](https://docs.mesg.com/guide/installation.html).

```bash
bash <(curl -fsSL https://mesg.com/install)
```

## Deploy the service

To deploy this service, go to [this service page](https://marketplace.mesg.com/services/{{.Sid}}) on the [MESG Marketplace](https://marketplace.mesg.com) and click the button "get/buy this service".

# Definitions

{{if .Events}}# Events
{{range $key, $event := .Events}}
## {{or $event.Name $key}}

Event key: `{{$key}}`

{{$event.Description}}

{{if $event.Data}}| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
{{range $dataKey, $data := $event.Data}}| **{{or $data.Name $dataKey}}** | `{{$dataKey}}` | `{{$data.Type}}` | {{if $data.Optional}}**`optional`** {{end}}{{$data.Description}} |
{{end}}{{end}}{{end}}{{end}}
{{if .Tasks}}# Tasks
{{range $key, $task := .Tasks}}
## {{or $task.Name $key}}

Task key: `{{$key}}`

{{$task.Description}}

{{if $task.Inputs}}### Inputs

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
{{range $inputKey, $input := $task.Inputs}}| **{{or $input.Name $inputKey}}** | `{{$inputKey}}` | `{{$input.Type}}` | {{if $input.Optional}}**`optional`** {{end}}{{$input.Description}} |
{{end}}{{end}}
{{if $task.Outputs}}### Outputs

{{range $outputKey, $output := $task.Outputs}}#### {{or $output.Name $outputKey}}

Output key: `{{$outputKey}}`

{{$output.Description}}

| **Name** | **Key** | **Type** | **Description** |
| --- | --- | --- | --- |
{{range $outputKey, $output := $output.Data}}| **{{or $output.Name $outputKey}}** | `{{$outputKey}}` | `{{$output.Type}}` | {{if $output.Optional}}**`optional`** {{end}}{{$output.Description}} |
{{end}}
{{end}}{{end}}{{end}}{{end}}
