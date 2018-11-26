# Test

While developing, you will want to listen to the different events that your Service might emit and also execute the different tasks that your Service provides.

## Run your Service in dev mode

The dev mode is a command that lets you start to monitor your Service. It lets you:
- start your Service
- log events from your Service
- log results from the tasks of your Service
- display the logs of your Service

```bash
mesg-core service dev ./PATH_TO_SERVICE_FOLDER
```

If you don't specify the path to the service folder, the command searches in the current folder for the `mesg.yml` file.

[More details here](../../cli/mesg-core_service_dev.md)

## Execute a task

With the `service dev` command your Service is up and running, but you will also need to execute specific tasks.

In order to do that, you need to get the generated `SERVICE_ID` from the `service dev` command and use it in the following command:

```bash
mesg-core service execute --task taskX --json TASK_INPUTS_JSON_FILE SERVICE_ID
```

The file for the inputs should be a `json` with a map of all the inputs that your task needs. For example:

```javascript
{
    "inputX": "...",
    "inputY": "..."
}
```

[More details here](../../cli/mesg-core_service_execute.md)

::: tip Get Help
You need help ? Check out the <a href="https://forum.mesg.com" target="_blank">MESG Forum</a>.
