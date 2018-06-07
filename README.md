<p align="center">
  <img src="https://cdn.rawgit.com/mesg-foundation/core/149-update-readme/logo.svg" alt="MESG Core" height="120">
</p>
<h2 align="center">
  <a href="https://mesg.tech/">Website</a> - 
  <a href="https://docs.mesg.tech/">Docs</a> - 
  <a href="https://medium.com/mesg">Blog</a> - 
  <a href="https://discordapp.com/invite/5tVTHJC">Discord</a>
</h2>

<p align="center">
  <a href="https://github.com/mesg-foundation/core"><img src="https://img.shields.io/circleci/project/github/mesg-foundation/core.svg" alt="CircleCI"></a>
  <a href="https://hub.docker.com/r/mesg/core/"><img src="https://img.shields.io/docker/pulls/mesg/core.svg" alt="Docker Pulls"></a>
  <a href="https://codeclimate.com/github/mesg-foundation/core/maintainability"><img src="https://api.codeclimate.com/v1/badges/86ad77f7c13cde40807e/maintainability" alt="Maintainability"></a>
  <a href="https://codecov.io/gh/mesg-foundation/core"><img src="https://codecov.io/gh/mesg-foundation/core/branch/dev/graph/badge.svg" alt="codecov"></a>
</p>

MESG is a platform for the creation of efficient and easy-to-maintain applications that connect any and all technologies.

# Quick Start Guide

### **1. Download the CLI**

First, download the CLI so you're able to interact with Core. You can either download the binaries directly from the [release page](https://github.com/mesg-foundation/core/releases/latest) then rename it to `mesg-core` and install it your path, or you can follow the installation process for your system:

### **2. Run MESG Core**

MESG needs to have a daemon running to process all the different commands that you might need to execute. In order to start the daemon you can run:

```text
mesg-core daemon start
```

### **3. Deploy a service**

Next step is to deploy the service that your application will need. You can [create your own service](https://docs.mesg.tech/service/what-is-a-service), but for now, let's just use an existing one and deploy it.

```text
mesg-core deploy https://github.com/mesg-foundation/service-webhook
```

Let's deploy another one.

```text
mesg-core deploy https://github.com/mesg-foundation/service-invite-discord
```

Every time you deploy a service, the console will display the ID for the service you've just deployed.

### **4. Connect the services**

Now, let's connect these services and create our application that will send you an email with an invitation to the MESG Discord every time you call the webhook.

```text
npm init && npm install --save mesg
```

Now create an `index.js` file and add the following code:

```javascript
const MESG = require('mesg/application')

const webhook    = '__ID_SERVICE_WEBHOOK__'
const invitation = '__ID_SERVICE_INVITATION_DISCORD__'
const email      = '__YOUR_EMAIL_HERE__'

MESG.ListenEvent({ serviceID: webhook, eventFilter: 'request' })
  .on('data', data => MESG.ExecuteTask({
    serviceID: invitation,
    taskKey: 'invite',
    taskData: JSON.stringify({ email })
  }, console.log))
```

Don't forget to replace the values `__ID_SERVICE_WEBHOOK__`, `__ID_SERVICE_INVITATION_DISCORD__` and `__YOUR_EMAIL_HERE__`.

### **5. Start the application**

Start your application now like any node application:

```javascript
npm start
```

### **6. Test the application**

Now we need to call the webhook in order to trigger the email, so let's do that with a curl command:

```text
curl -XPOST http://localhost:3000/webhook
```

You should now have an email in your inbox with your precious invitation to our Discord.


# Build from source

## Download source

```bash
mkdir -p $GOPATH/src/github.com/mesg-foundation/core
cd $GOPATH/src/github.com/mesg-foundation/core
git clone https://github.com/mesg-foundation/core.git ./
```

## Install dependencies

```bash
go get -v -t -u ./...
```

## Run all tests with code coverage

```bash
env CORE.IMAGE=mesg/core:local go test -cover -v ./...
```

If you use Visual code you can add the following settings (Preference > Settings)
```json
"go.testEnvFile": "${workspaceRoot}/testenv"
```

## Build MESG Core and start it

```bash
./dev-core
```

## Build CLI and start it

```bash
./dev-cli
```

## Install debugger on OS X

```bash
xcode-select --install
go get -u github.com/derekparker/delve/cmd/dlv
```
If the debugger still doesn't work, try the following:
```bash
cd $GOPATH/src/github.com/derekparker/delve
make install
```

[Source](https://github.com/derekparker/delve/blob/master/Documentation/installation/osx/install.md)

# Issues

For issues concerning application or service development, please read the [docs](https://docs.mesg.tech/) or ask us directly on the [Discord](https://discordapp.com/invite/5tVTHJC) channels #application or #service.

For questions or suggestions regarding a new feature of Core, please contact us on [Discord](https://discordapp.com/invite/5tVTHJC) channel #core.

To report a bug, please [check for existing issues, then create a new issue on this repository](https://github.com/mesg-foundation/core/issues).

# Contribution

For Services and Applications contribution, we have an [curated list of Awesome Services and Applications](https://github.com/mesg-foundation/awesome) that you should participate in.

For contribution to MESG's Core, please contact us on [Discord](https://discordapp.com/invite/5tVTHJC) channel #core. We would love to include you in the development process.

## Service

MESG depends heavily on services. These services are automatically built and ran inside Docker. You can connect anything you want, as long as it can run inside Docker \(as long as it can run on a computer\). If you need more details about how to connect dependencies to your service [check out the documentation](https://docs.mesg.tech/service/dockerize-the-service).

A service needs to implement two types of communications:receiving tasks and submitting events.

#### Receiving Tasks

Tasks are designed to receive information from Core and the Application that you run. Tasks can have multiple parameters as inputs and multiple outputs with varying data. You can visualize a task as a simple function that can return any kind of object.

You could have a task that takes a name as an input, and shows `success` as an output. This task factors the type of name with its probability like `{ "type": "female", "proabiliy": 92.34% }` but could also have an `error` output with a type of error like this `{ "message": "This doesn't looks like a name" }`.

Check out the documentation for more information on how info how to create [tasks](https://docs.mesg.tech/service/listen-for-tasks).

#### Submitting Events

Let's say you are working with a webserver. An event could be when there is a request with data in the payload, or it could be different events for each of the different routes of your API, or, in a blockchain context, it could be when a smart contract emits an event.

For more info how to create your events, visit the [Emit an Event](https://docs.mesg.tech/service/emit-an-event) page.

## Architecture

\[\[ TODO: Add a nice graphic with the Application, the core and the services with the communication \]\]

## Examples

You can find a list of different examples and services that you can re-use [here](https://github.com/mesg-foundation/awesome)

## Roadmap
