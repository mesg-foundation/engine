## mesg-core service delete

Delete one or many services

### Synopsis

Delete one or many services

```
mesg-core service delete [flags]
```

### Examples

```
mesg-core service delete SERVICE [SERVICE...]
mesg-core service delete --all
```

### Options

```
      --all         Delete all services
  -h, --help        help for delete
      --keep-data   Do not delete services' persistent data
  -y, --yes         Automatic "yes" to all prompts and run non-interactively
```

### Options inherited from parent commands

```
      --no-color     disable colorized output
      --no-spinner   disable spinners
```

### SEE ALSO

* [mesg-core service](mesg-core_service.md)	 - Manage services

