## mesg-cli service start

Start a service

### Synopsis

Start a service from the published available services. You have to provide a stake value and duration.

```
mesg-cli service start SERVICE_ID [flags]
```

### Examples

```
mesg-cli service start SERVICE_ID
mesg-cli service start SERVICE_ID --stake STAKE --duration DURATION  --account ACCOUNT --confirm
```

### Options

```
  -a, --account string   Account to use
  -c, --confirm          Confirm
  -d, --duration int     The duration you will be running this service
  -h, --help             help for start
  -s, --stake float      The amount to stake
```

### SEE ALSO

* [mesg-cli service](mesg-cli_service.md)	 - Manage your services

