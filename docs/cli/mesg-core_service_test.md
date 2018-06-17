## mesg-core service test

Test a service

### Synopsis

Test a service by listening to events or calling tasks.

See more detail on the [Test page from the documentation](https://docs.mesg.tech/service/test.html)

```
mesg-core service test [flags]
```

### Examples

```
mesg-core service test
mesg-core service test ./SERVICE_FOLDER
mesg-core service test --event-filter EVENT_NAME
mesg-core service test --task TASK_NAME --data ./PATH_TO_DATA_FILE.json
mesg-core service test --task-filter TASK_NAME --output-filter OUTPUT_NAME
mesg-core service test --serviceID SERVICE_ID --keep-alive
```

### Options

```
  -d, --data string            Path to the file containing the data required to run the task
  -e, --event-filter string    Only log the data of the given event (default "*")
  -f, --full-logs              Display logs from service and its dependencies
  -h, --help                   help for test
      --keep-alive             Do not stop the service at the end of this command
  -o, --output-filter string   Only log the data of the given output of a task result. If set, you also need to set the task in --task-filter
  -s, --serviceID string       ID of a previously deployed service
  -t, --task string            Run the given task
  -r, --task-filter string     Only log the result of the given task
```

### SEE ALSO

* [mesg-core service](mesg-core_service.md)	 - Manage your services

