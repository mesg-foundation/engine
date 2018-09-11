## mesg-core service stop

Stop a service

### Synopsis

Stop a service.

**WARNING:** If you stop a service with your stake duration still ongoing, you may lost your stake.
You will **NOT** get your stake back immediately. You will get your remaining stake only after a delay.
To have more explanation, see the page [stake explanation from the documentation]().

```
mesg-core service stop SERVICE [flags]
```

### Examples

```
mesg-core service stop SERVICE
```

### Options

```
  -h, --help   help for stop
```

### Options inherited from parent commands

```
      --no-color     disable colorized output
      --no-spinner   disable spinners
```

### SEE ALSO

* [mesg-core service](mesg-core_service.md)	 - Manage your services

