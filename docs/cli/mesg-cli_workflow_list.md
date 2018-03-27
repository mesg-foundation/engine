## mesg-cli workflow list

List all deployed workflows

### Synopsis

List all workflows deployed on the Network.

Optionally, you can filter the workflows deployed by a specific account.

This command will return basic information. To have more details, see the [detail command](mesg-cli_workflow_detail.md).

```
mesg-cli workflow list [flags]
```

### Examples

```
mesg-cli workflow list
mesg-cli workflow list --account ACCOUNT
```

### Options

```
  -a, --account string   Filter workflows by a specific account
  -h, --help             help for list
```

### SEE ALSO

* [mesg-cli workflow](mesg-cli_workflow.md)	 - Manage your workflows

