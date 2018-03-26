## mesg-cli service stop

Stop a service

### Synopsis

Stop a service.

**WARNING:** If you stop a service with your stake duration still ongoing, you may lost your stake.
You will **NOT** get your stake back immediately. You will get your remaining stake only after a delay.
To have more explanation, see the page [stake explanation from the documentation](https://docs.mesg.tech/service/run/).
	

```
mesg-cli service stop SERVICE_ID [flags]
```

### Examples

```
mesg-cli service stop SERVICE_ID
mesg-cli service stop SERVICE_ID --account ACCOUNT --confirm
```

### Options

```
  -a, --account string   Account you want to use
  -c, --confirm          Confirm
  -h, --help             help for stop
```

### SEE ALSO

* [mesg-cli service](mesg-cli_service.md)	 - Manage your services

