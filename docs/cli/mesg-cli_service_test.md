## mesg-cli service test

Start and test the service

### Synopsis

Test the interactions with the service, listening to events and calling tasks.

```
mesg-cli service test SERVICE_FILE [flags]
```

### Examples

```
mesg-cli service test service.yml
```

### Options

```
  -d, --data string    File with the data required to run a specific task
  -e, --event string   Event filter, will only log those events
  -h, --help           help for test
      --keep-alive     Let the service run event after the end of the test
  -t, --task string    Run a specific task
```

### SEE ALSO

* [mesg-cli service](mesg-cli_service.md)	 - Manage the services you created

