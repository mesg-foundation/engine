## mesg-cli account export

Export account details in order to be able to re-import it with the import command

### Synopsis

Export account details in order to be able to re-import it with the import command

```
mesg-cli account export [flags]
```

### Examples

```
mesg-cli account export --account AccountX
```

### Options

```
  -a, --account string        Account you want to use
  -h, --help                  help for export
      --new-password string   New password for the account you export
      --password string       Current password for the account you export
  -p, --path string           Path of the file where your account will be exported (default "./export")
```

### SEE ALSO

* [mesg-cli account](mesg-cli_account.md)	 - Manage your MESG accounts

