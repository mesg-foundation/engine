## mesg-core service logs

Show logs of a service

### Synopsis

Show logs of a service

```
mesg-core service logs [flags]
```

### Examples

```
mesg-core service logs SERVICE
mesg-core service logs SERVICE --dependencies DEPENDENCY_NAME,DEPENDENCY_NAME,...
```

### Options

```
  -d, --dependencies stringArray   Name of the dependency to show the logs from
  -h, --help                       help for logs
```

### Options inherited from parent commands

```
      --no-color     disable colorized output
      --no-spinner   disable spinners
```

### SEE ALSO

* [mesg-core service](mesg-core_service.md)	 - Manage services

