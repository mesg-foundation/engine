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
mesg-core service test --event EVENT_NAME
mesg-core service test --task TASK_NAME --data ./PATH_TO_DATA_FILE.yml
mesg-core service test --keep-alive
```

### Options

```
  -d, --data string    Path to the file containing the data required to run the task
  -e, --event string   Only log a specific event (default "*")
  -h, --help           help for test
  -t, --task string    Run a specific task
```

### SEE ALSO

* [mesg-core service](mesg-core_service.md)	 - Manage your services

