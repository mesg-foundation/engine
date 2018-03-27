## mesg-cli workflow resume

Resume a previously paused workflow

### Synopsis

Resume a previously paused workflow.

To pause a workflow, see the [pause command](mesg-cli_workflow_pause.md)

```
mesg-cli workflow resume WORKFLOW_ID [flags]
```

### Examples

```
mesg-cli workflow resume WORKFLOW_ID
mesg-cli workflow resume WORKFLOW_ID --account ACCOUNT_ID --confirm
```

### Options

```
  -a, --account string   Account to use
  -c, --confirm          Confirm
  -h, --help             help for resume
```

### SEE ALSO

* [mesg-cli workflow](mesg-cli_workflow.md)	 - Manage your workflows

