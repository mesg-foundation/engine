## mesg-cli marketplace service pause

Pause a service

### Synopsis

Pause a service. The user will not get its stake back but it will also not lost it. Should always pause services before quitting the CLI, otherwise the user may loss its stake.

```
mesg-cli marketplace service pause SERVICE [flags]
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

* [mesg-cli marketplace service](mesg-cli_marketplace_service.md)	 - Manage services from the marketplace

