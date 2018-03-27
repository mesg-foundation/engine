## mesg-cli service list

List all published services

### Synopsis

This command returns all published services with basic information.
Optionally, you can filter the services published by a specific developer:
To have more details, see the [detail command](mesg-cli_service_detail.md).

```
mesg-cli service list [flags]
```

### Examples

```
mesg-cli service list
mesg-cli service list --account ACCOUNT
```

### Options

```
  -a, --account string   Account to use
  -h, --help             help for list
```

### SEE ALSO

* [mesg-cli service](mesg-cli_service.md)	 - Manage your services

