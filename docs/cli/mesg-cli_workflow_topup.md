## mesg-cli workflow topup

Top-up a workflow

### Synopsis

Top-up a workflow.
Add more token to a existing workflow.

```
mesg-cli workflow topup WORKFLOW_ID [flags]
```

### Examples

```
mesg-cli workflow topup WORKFLOW_ID
mesg-cli workflow topup WORKFLOW_ID --amount AMOUNT --account ACCOUNT_ID --confirm
```

### Options

```
  -a, --account string   Account to use
      --amount string    Amount of MESG
  -c, --confirm          Confirm
  -h, --help             help for topup
```

### SEE ALSO

* [mesg-cli workflow](mesg-cli_workflow.md)	 - Manage your workflows

