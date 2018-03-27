## mesg-cli service resume

Resume a service

### Synopsis

Resume a previously paused service.

To pause a service, see the [pause command](mesg-cli_service_pause.md)

```
mesg-cli service resume SERVICE_ID [flags]
```

### Examples

```
mesg-cli service resume SERVICE_ID --account ACCOUNT --confirm
```

### Options

```
  -a, --account string   Account to use
  -c, --confirm          Confirm
  -h, --help             help for resume
```

### SEE ALSO

* [mesg-cli service](mesg-cli_service.md)	 - Manage your services

