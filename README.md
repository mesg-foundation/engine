![MESG core](https://camo.githubusercontent.com/2b99bc67988c4793a28d04cc471e10948da15ff3/68747470733a2f2f63646e2e646973636f72646170702e636f6d2f6174746163686d656e74732f3435303131353235303838313033363334382f3435313730353138363835363037353236352f4d4553472d4769746875622d62616e332d30322e6a7067)

[Website](https://mesg.tech/) - [Docs](https://docs.mesg.tech/) - [Chat](https://discordapp.com/invite/SaZ5HcE) - [Blog](https://medium.com/mesg)


[![CircleCI](https://img.shields.io/circleci/project/github/mesg-foundation/core.svg)](https://github.com/mesg-foundation/core)
[![Docker Pulls](https://img.shields.io/docker/pulls/mesg/daemon.svg)](https://hub.docker.com/r/mesg/daemon/)
[![Maintainability](https://api.codeclimate.com/v1/badges/86ad77f7c13cde40807e/maintainability)](https://codeclimate.com/github/mesg-foundation/core/maintainability)
[![codecov](https://codecov.io/gh/mesg-foundation/core/branch/dev/graph/badge.svg)](https://codecov.io/gh/mesg-foundation/core)




**MESG is a platform to help you create efficient and easy to maintain applications that connects any technologie**. You can create your application that listen to an event on a Blockchain and execute a task on a web server. The technology doesn't matter, as long as it can send and/or receive some data.

The idea with MESG is to create different services that connects to a specific technology and/or a specific feature and expose 2 different things:
- Events
- Tasks

When these services are created you can then create your business logic with the event driven paradygme and connect any **event** from a service to a **task** of another service.

# Contents

- [Quickstart](quickstart)
- [Service](service)
  - [Receiving Task](receiving-task)
  - [Submitting Event](submitting-event)
- [Architecture](architecture)
- [Examples](examples)
- [Roadmap](roadmap)
- [Join us](join-us)

# Quickstart

### 1 - Download the CLI

First thing you need to do is to download the CLI to be able to interact with the MESG Core.
You can either download the binaries directly from the [release page](https://github.com/mesg-foundation/core/releases/latest) then rename it to `mesg-core` and install it your path or you can follow the details from [the documentation](https://docs.mesg.tech/start-here/installing-core).

One this is done open a new terminal and type `mesg-core` and you should have something similar to this.

[[ TODO: Insert screenshot of the command line ]]

### 2 - Run the MESG Daemon

MESG needs to have a daemon running to process all the different commands that you might need to execute. In order to start the daemon you can run:
```
mesg-core daemon start
```

### 3 - Deploy a service

Next step is to deploy the service that your application will need. You can [create your own service](https://docs.mesg.tech/service/what-is-a-service) but for now let's just use an existing one and deploy it.

```
mesg-core deploy https://github.com/mesg-foundation/service-webhook
```

Let's deploy another one.

```
mesg-core deploy https://github.com/mesg-foundation/service-invite-discord
```

Every time you deploy, the console will display you the ID for the service you've just deployed.

### 4 - Connect the services

Now let's connect these services and create our application that will send you an email with an invitation to the MESG Discord every time you call the webhook.

```
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

### 5 - Start the application

Start now your application like any node application:
```javascript
npm start
```

### 6 - Test the application

Now we need to call the webhook in order to trigger the email so let's do that with a curl command.
```
curl -XPOST http://localhost:3000/webhook
```
You should now have an email in your inbox with you precious invitation to our Discord.

# Service

MESG depends heavily on services. These services are automatically build and run inside Docker. You can connect anything you want as long as it can run inside Docker (so as long as it can run in a computer). If you need more details about how to connect dependencies to your service [checkout the documentation](https://docs.mesg.tech/service/dockerize-the-service).

Your service needs to implement two types of communications:

## Receiving Task

Task are designed to receive informations from the MESG core and the Application that you run. Tasks can have multiple parameters as inputs and muliple outputs with multiple data. Visualize a task as a simple function that can return any kind of object.

You could have a task that take as input a name and as output `success` and this task the genre of this name with the probability like `{ "genre": "female", "proabiliy": 92.34% }` but could also have an `error` output with the type of error `{ "message": "This doesn't looks like a name" }`.

More info how to create your [tasks in the documentation](https://docs.mesg.tech/service/listen-for-tasks).

## Submitting Event

Events are data that your service will emit in real time. Let's say you are doing a webserver. One event could be when there is a request with the data in the payload or different events for the different routes of your api or in a blockchain world when a smart contract emits an event.

More info how to create your [events in the documentation](https://docs.mesg.tech/service/emit-an-event)

# Architecture

[[ TODO: Add a nice graphic with the Application, the core and the services with the communication ]]

# Examples

You can find a list of different examples and services that you can re-use [here](https://github.com/mesg-foundation/awesome)

# Roadmap


# Join us


## Build from source

### Download source

```bash
mkdir -p $GOPATH/src/github.com/mesg-foundation/core
cd $GOPATH/src/github.com/mesg-foundation/core
git clone https://github.com/mesg-foundation/core.git ./
```

### Install dependencies

```bash
go get -v -t ./...
```

### Run all tests with code coverage

```bash
env DAEMON.IMAGE=mesg/daemon:local go test -cover -v ./...
```

If you use Visual code you can add the following settings (Preference > Settings)
```json
"go.testEnvFile": "${workspaceRoot}/testenv"
```

### Build docker image

```bash
docker build -t mesg-daemon .
```

### Install debugger on OS X

```bash
xcode-select --install
go get -u github.com/derekparker/delve/cmd/dlv
cd $GOPATH/src/github.com/derekparker/delve
make install
```

[Source](https://github.com/derekparker/delve/blob/master/Documentation/installation/osx/install.md)

