## mesg-core service dev

Run your service in development mode

### Synopsis

Run your service in development mode

```
mesg-core service dev [flags]
```

### Examples

```
mesg-core service dev PATH
```

### Options

```
  -e, --event-filter string    Only log the data of the given event (default "*")
  -h, --help                   help for dev
  -o, --output-filter string   Only log the data of the given output of a task result. If set, you also need to set the task in --task-filter
  -t, --task-filter string     Only log the result of the given task
```

### Options inherited from parent commands

```
      --no-color     disable colorized output
      --no-spinner   disable spinners
```

### SEE ALSO

* [mesg-core service](mesg-core_service.md)	 - Manage services

