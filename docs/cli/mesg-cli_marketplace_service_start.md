## mesg-cli marketplace service start

Start a service

### Synopsis

Start a service from the publicly available services. The user have to provide a stake value and duration.

```
mesg-cli marketplace service start SERVICE [flags]
```

### Examples

```
mesg-cli marketplace service start --stake 100 --duration 10 ethereum
```

### Options

```
  -a, --account string   Account you want to use
  -c, --confirm          Confirm
  -d, --duration int     The amount of time you will be running this/those service(s) for (in hours)
  -h, --help             help for start
  -s, --stake float      The number of MESG to put on stake
```

### SEE ALSO

* [mesg-cli marketplace service](mesg-cli_marketplace_service.md)	 - Manage services from the marketplace

