## mesg-core service logs

Show the logs of a service

### Synopsis

Show the logs of a service

```
mesg-core service logs [flags]
```

### Examples

```
mesg-core service logs SERVICE
mesg-core service logs SERVICE --dependency DEPENDENCY_NAME
```

### Options

```
  -d, --dependency string   Name of the dependency to show the logs from (default "*")
  -h, --help                help for logs
```

### Options inherited from parent commands

```
      --no-color     disable colorized output
      --no-spinner   disable spinners
```

### SEE ALSO

* [mesg-core service](mesg-core_service.md)	 - Manage your services

