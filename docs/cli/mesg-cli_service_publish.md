## mesg-cli service publish

Publish a new service

### Synopsis

Deploy a service to the Network from a given service file. Validate it first. The user will need to provide an account and the password of the account.

```
mesg-cli service publish SERVICE_FILE [flags]
```

### Examples

```
mesg-cli service publish service.yml
```

### Options

```
  -a, --account string   Account you want to use
  -c, --confirm          Confirm
  -h, --help             help for publish
```

### SEE ALSO

* [mesg-cli service](mesg-cli_service.md)	 - Manage the services you created

