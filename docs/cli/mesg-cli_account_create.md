## mesg-cli account create

Create a new account

### Synopsis

This method creates a new account secured by a password. We strongly advise to use long randomized password.

**Warning:** Backup your password in a safe place. You will not be able to use the account if you lost the password.

You should also [export your account](account/export.md) to a safe place to prevent losing access to your workflows, services and tokens.

```
mesg-cli account create [flags]
```

### Examples

```
mesg-cli account create
mesg-cli account create --password QWERTY
```

### Options

```
  -h, --help              help for create
  -p, --password string   Password of the account
```

### SEE ALSO

* [mesg-cli account](mesg-cli_account.md)	 - Manage your accounts

