## mesg-cli account create

Create a new account

### Synopsis

This method creates a new account secured by your password. We strongly advise to use long randomized password.
Warning: Backup your password in a safe place. You will not be able to use the account if you lost the password.
You should also export your account to a safe place to prevent losing access to your workflows, services and tokens. See account export method.

```
mesg-cli account create [flags]
```

### Examples

```
mesg-cli account create --name ACCOUNT_NAME --password ACCOUNT_PASSWORD
```

### Options

```
  -h, --help              help for create
  -n, --name string       Name of the account
  -p, --password string   Password for the account
```

### SEE ALSO

* [mesg-cli account](mesg-cli_account.md)	 - Manage your MESG accounts

