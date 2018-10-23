# Installation

## Automatic installation

Run the following command in a console:

```bash
bash <(curl -fsSL https://mesg.com/install)
```

## Manual installation

* Download and install [Docker CE](https://www.docker.com/community-edition)
* Download the binary from our [release page on GitHub](https://github.com/mesg-foundation/core/releases)
* Rename the binary to `mesg-core`
* Give it the execution permission: `chmod +x mesg-core`
* Move it to your local bin folder: `mv ./mesg-core /usr/local/bin/mesg-core`
* Clone the core repo to get system services: `git clone https://github.com/mesg-foundation/core.git /tmp/mesg-core`
* Create system services folder under your mesg path: `mkdir -p ~/.mesg/systemservices`
* Copy system services from core repo: `cp -a /tmp/mesg-core/systemservices/-sources/. ~/.mesg/systemservices`
* Start MESG Core with the command: `mesg-core start`


::: tip Get Help
You need help ? Check out the <a href="https://forum.mesg.com" target="_blank">MESG Forum</a>.