## mesg-core service init

Initialize a service

### Synopsis

Initialize a service by creating a mesg.yml and Dockerfile in a dedicated directory.
	
To get more information, see the page [service file from the documentation](https://docs.mesg.com/guide/service/service-file.html)

```
mesg-core service init [flags]
```

### Examples

```
mesg-core service init
mesg-core service init --name NAME --description DESCRIPTION
mesg-core service init --current
```

### Options

```
  -d, --description string   Service description
      --dir string           Create the service in the direcotry
  -h, --help                 help for init
  -n, --name string          Service name
  -t, --template string      Specify the template URL to use
```

### Options inherited from parent commands

```
      --no-color     disable colorized output
      --no-spinner   disable spinners
```

### SEE ALSO

* [mesg-core service](mesg-core_service.md)	 - Manage your services

