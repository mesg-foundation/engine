## mesg-cli account export

Export an account

### Synopsis

This method creates a file containing the information about your account.
The private key of your account is encrypted with the password you choose.

**Warning:** This method does **NOT** export your password. You have to manage your password yourself.

You can import the backup file on any other MESG Application with the [import method](mesg-cli_account_import.md).

```
mesg-cli account export [flags]
```

### Examples

```
mesg-cli account export
mesg-cli account export --account ACCOUNT --password PASSWORD --new-password PASSWORD --path ./PATH_TO_BACKUP_FILE
```

### Options

```
  -a, --account string        Account to use
  -h, --help                  help for export
      --new-password string   New password of the exported account
      --password string       Current password of the account to export
  -p, --path string           Path of the file your account will be exported in (default "./export")
```

### SEE ALSO

* [mesg-cli account](mesg-cli_account.md)	 - Manage your accounts

