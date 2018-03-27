## mesg-cli account delete

Delete an account

### Synopsis

This method deletes an account.

**Warning:** If you didn't previously [export this account](mesg-cli_account_export.md), you will lost it **forever.**

```
mesg-cli account delete [flags]
```

### Examples

```
mesg-cli service delete
mesg-cli service delete --account 0x0000000000000000000000000000000000000000 --confirm
```

### Options

```
  -a, --account string   Account to use
  -c, --confirm          Confirm
  -h, --help             help for delete
```

### SEE ALSO

* [mesg-cli account](mesg-cli_account.md)	 - Manage your accounts

