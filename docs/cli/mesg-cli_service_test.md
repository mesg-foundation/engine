## mesg-cli service test

Test a service

### Synopsis

Test a service by listening to events and calling tasks.

See more detail on the [Test page from the documentation](https://docs.mesg.tech/service/develop/test.html)

```
mesg-cli service test [flags]
```

### Examples

```
mesg-cli service test
mesg-cli service test ./SERVICE_FOLDER
```

### Options

```
  -d, --data string    Path to the file containing the data required to run the task
  -e, --event string   Only log a specific event
  -h, --help           help for test
      --keep-alive     Leave the service runs after the end of the test
  -t, --task string    Run a task
```

### SEE ALSO

* [mesg-cli service](mesg-cli_service.md)	 - Manage your services

