## mesg-cli workflow test

Test a workflow

### Synopsis

Test a workflow locally

To get more information, see the [test page from the documentation](https://docs.mesg.tech/workflow/test.html)

```
mesg-cli workflow test ./PATH_TO_WORKFLOW_FILE [flags]
```

### Examples

```
mesg-cli workflow test ./PATH_TO_WORKFLOW_FILE.yml
mesg-cli workflow test ./PATH_TO_WORKFLOW_FILE.yml --live
mesg-cli workflow test ./PATH_TO_WORKFLOW_FILE.yml --task TASK_ID --live
mesg-cli workflow test ./PATH_TO_WORKFLOW_FILE.yml --task TASK_ID --event ./PATH_TO_EVENT_DATA_FILE.yml
mesg-cli workflow test ./PATH_TO_WORKFLOW_FILE.yml --live --keep-alive
```

### Options

```
  -e, --event string   Path to the event data file
  -h, --help           help for test
  -k, --keep-alive     Keep the services alive (re-run without the option to stop them)
  -l, --live           Use live events
  -t, --task string    Run the test on a specific task of the workflow
```

### SEE ALSO

* [mesg-cli workflow](mesg-cli_workflow.md)	 - Manage your workflows

