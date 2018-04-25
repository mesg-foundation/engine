## mesg-cli service test

Test a service

### Synopsis

Test a service by listening to events or calling tasks.

See more detail on the [Test page from the documentation](https://docs.mesg.tech/service/test.html)

```
mesg-cli service test [flags]
```

### Examples

```
mesg-cli service test
mesg-cli service test ./SERVICE_FOLDER
mesg-cli service test --event EVENT_NAME
mesg-cli service test --task TASK_NAME --data ./PATH_TO_DATA_FILE.yml
mesg-cli service test --keep-alive
```

### Options

```
  -e, --event string   Only log a specific event (default "*")
  -h, --help           help for test
```

### SEE ALSO

* [mesg-cli service](mesg-cli_service.md)	 - Manage your services

