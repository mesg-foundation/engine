## mesg-core service init

Initialize a service

### Synopsis

Initialize a service by creating a mesg.yml and Dockerfile in a dedicated folder.
	
To get more information, see the page [service file from the documentation](https://docs.mesg.com/service/service-file.html)

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
  -c, --current              Create the service in the current path
  -d, --description string   Description
  -h, --help                 help for init
  -n, --name string          Name
  -t, --template string      Specify the template URL to use
```

### SEE ALSO

* [mesg-core service](mesg-core_service.md)	 - Manage your services

