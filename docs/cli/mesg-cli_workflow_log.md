## mesg-cli workflow log

Log executions of a workflow

### Synopsis

Log executions of a workflow

```
mesg-cli workflow log WORKFLOW_ID [flags]
```

### Examples

```
mesg-cli workflow log WORKFLOW_ID
	mesg-cli workflow log WORKFLOW_ID --account ACCOUNT_ID
mesg-cli workflow log WORKFLOW_ID --execution EXECUTION_ID
mesg-cli workflow log WORKFLOW_ID --task TASK_ID
mesg-cli workflow log WORKFLOW_ID --from DATE --to DATE
```

### Options

```
  -a, --account string     Account to use
  -e, --execution string   Log a specific execution
      --from string        Log from date
  -h, --help               help for log
  -t, --task string        Log a specific task
      --to string          Log to date
```

### SEE ALSO

* [mesg-cli workflow](mesg-cli_workflow.md)	 - Manage your workflows

