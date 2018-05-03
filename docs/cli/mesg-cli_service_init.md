## mesg-cli service init

Initialize a service

### Synopsis

Initialize a service by creating a mesg.yml and Dockerfile in a dedicated folder.
	
To get more information, see the page [service file from the documentation](https://docs.mesg.tech/service/service-file.html)

```
mesg-cli service init [flags]
```

### Examples

```
mesg-cli service init
mesg-cli service init --name NAME --description DESCRIPTION --visibility ALL --publish ALL
```

### Options

```
  -c, --current              Create the service in the current path
  -d, --description string   Description
  -h, --help                 help for init
  -n, --name string          Name
  -p, --publish string       Publish
  -v, --visibility string    Visibility
```

### SEE ALSO

* [mesg-cli service](mesg-cli_service.md)	 - Manage your services

