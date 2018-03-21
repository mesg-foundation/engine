## mesg-cli service pause

Pause a service

### Synopsis

Pause a service. The user will not get its stake back but it will also not lost it. Should always pause services before quitting the CLI, otherwise the user may loss its stake.

```
mesg-cli service pause SERVICE [flags]
```

### Examples

```
mesg-cli marketplace service pause ethereum
```

### Options

```
  -a, --account string   Account you want to use
  -c, --confirm          Confirm
  -h, --help             help for pause
```

### SEE ALSO

* [mesg-cli service](mesg-cli_service.md)	 - Manage the services you created

