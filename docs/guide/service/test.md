# Test

Before deploying the Service, you'll want to test it to ensure that everything is working properly.

To ensure that the Application is able to start your Service and receive an event from it, execute the following method:

```bash
mesg-core service test ./PATH_TO_SERVICE_FOLDER
```

If you don't specify the path to the service folder, the command searches in the current folder for the `mesg.yml` file.

## Listen to a specific event

To only listen to one specific event, you can specify the name of the event with the `--event-filter` flag:

```bash
mesg-core service test --event-filter myServiceEventName
```

## Run a task

To test a task from your Service, run:

```bash
mesg-core service test --task myServiceTaskName
```

If your task requires inputs you will need to specify the file that contains all the input values in `json` format.

```bash
mesg-core service test --task myServiceTaskName --data ./PATH_TO_DATA_FILE.json
```

The file should be a `json` with a format similar to the following:

```javascript
{
    "inputX": "...",
    "inputY": "..."
}
```

## Listen to a specific task

To listen to the result of a specific task, you can use the flag `--task-filter`

```bash
mesg-core service test --task-filter myServiceTaskName
```

You can also listen to a specific output of the result of a task by using the flag `--task-filter` and `--output-filter`

```bash
mesg-core service test --task-filter myServiceTaskName --output-filter myServiceOutputName
```

## Keep it alive

All previous commands will stop your service upon quitting. If you want to leave your service alive, you can add the following flag to any command: `--keep-alive`. For example:

```bash
mesg-core service test --task myServiceTaskName --keep-alive
```

## Test an already deployed service

If you want to test a service that has already been deployed, you can pass its ID to the flag `--serviceID`

```bash
mesg-core service test --serviceID myServiceID
```



