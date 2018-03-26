## mesg-cli service pause

Pause a service

### Synopsis

Pause a service. The service should have been previously started.

You should always pause services before quitting the CLI, otherwise you may loss your stake.

You will **NOT** get your stake back with this command. The goal of this command is to give you an opportunity to stop running a service for a short period of time without losing your stake.

When a service is paused, the stake duration count is also paused.

```
mesg-cli service pause SERVICE_ID [flags]
```

### Examples

```
mesg-cli service pause SERVICE_ID
mesg-cli service pause SERVICE_ID --account ACCOUNT --confirm
```

### Options

```
  -a, --account string   Account you want to use
  -c, --confirm          Confirm
  -h, --help             help for pause
```

### SEE ALSO

* [mesg-cli service](mesg-cli_service.md)	 - Manage your services

