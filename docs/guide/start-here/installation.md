# Installation

## Automatic installation

Run the following command in a console:

```bash
bash <(curl -fsSL https://mesg.com/install)
```

## Manual installation

### Docker

* Download and install [Docker CE](https://www.docker.com/community-edition)
* Initialize Docker Swarm by running
```
docker swarm init
```
* If the error `Could not choose an IP address to advertise since this system has multiple addresses on interface eth0 (xxx.xxx.xxx.xxx and yyy.yyy.yyy.yyy)` is returned, run:
```
docker swarm init --advertise-addr xxx.xxx.xxx.xxx
```

### CLI

* Download the binary from our [release page on GitHub](https://github.com/mesg-foundation/core/releases)
* Rename the binary to `mesg-core`
* Give it the execution permission
```
chmod +x mesg-core
```
* Move it to your local bin folder
```
mv ./mesg-core /usr/local/bin/mesg-core
```
<!-- * Clone the core repo to get system services
```
git clone https://github.com/mesg-foundation/core.git /tmp/mesg-core
```
* Create system services folder under your mesg path
```
mkdir -p ~/.mesg/systemservices
```
* Copy system services from core repo
```
cp -a /tmp/mesg-core/systemservices/sources/. ~/.mesg/systemservices -->
```
* Start MESG Core with the command
```
mesg-core start
```

### Docker only

If you don't want to use the CLI, you can start the Core by executing the following commands.

* Download latest version
```
docker pull mesg/core:latest
```

* Create the MESG network
```
docker network create core -d overlay --label com.docker.stack.namespace=core
```

* Start MESG Core
```
docker service create --network core --env MESG_CORE_PATH=/mesg --mount source=/var/run/docker.sock,destination=/var/run/docker.sock,type=bind --mount source=$HOME/.mesg,destination=/mesg,type=bind --publish 50052:50052 --label com.docker.stack.namespace=core --name core mesg/core:latest
```

::: tip Get Help
You need help ? Check out the <a href="https://forum.mesg.com" target="_blank">MESG Forum</a>.