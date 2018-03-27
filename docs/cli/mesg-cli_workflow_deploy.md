## mesg-cli workflow deploy

Deploy a workflow

### Synopsis

Deploy a workflow on the Network.

To get more information, see the [deploy page from the documentation](https://docs.mesg.tech/workflow/deploy.html)

```
mesg-cli workflow deploy ./PATH_TO_WORKFLOW_FILE [flags]
```

### Examples

```
mesg-cli workflow deploy ./PATH_TO_WORKFLOW_FILE.yml
mesg-cli workflow deploy ./PATH_TO_WORKFLOW_FILE.yml --account ACCOUNT --amount AMOUNT --confirm
```

### Options

```
  -a, --account string   Account to use
      --amount string    Amount of MESG
  -c, --confirm          Confirm
  -h, --help             help for deploy
```

### SEE ALSO

* [mesg-cli workflow](mesg-cli_workflow.md)	 - Manage your workflows

