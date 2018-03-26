## mesg-cli account import

Import an account from a backup file

### Synopsis

This method imports a previously exported backup file of your account created with the [export method](account/export.md).

```
mesg-cli account import FILE [flags]
```

### Examples

```
mesg-cli account import ./export
```

### Options

```
  -h, --help                  help for import
      --new-password string   New password of the imported account
      --password string       Current password of the account to import
```

### SEE ALSO

* [mesg-cli account](mesg-cli_account.md)	 - Manage your accounts

