![MESG core](https://camo.githubusercontent.com/2b99bc67988c4793a28d04cc471e10948da15ff3/68747470733a2f2f63646e2e646973636f72646170702e636f6d2f6174746163686d656e74732f3435303131353235303838313033363334382f3435313730353138363835363037353236352f4d4553472d4769746875622d62616e332d30322e6a7067)

[Website](https://mesg.tech/) - [Docs](https://docs.mesg.tech/) - [Chat](https://discordapp.com/invite/SaZ5HcE) - [Blog](https://medium.com/mesg)


[![CircleCI](https://img.shields.io/circleci/project/github/mesg-foundation/core.svg)](https://github.com/mesg-foundation/core)
[![Docker Pulls](https://img.shields.io/docker/pulls/mesg/daemon.svg)](https://hub.docker.com/r/mesg/daemon/)
[![Maintainability](https://api.codeclimate.com/v1/badges/86ad77f7c13cde40807e/maintainability)](https://codeclimate.com/github/mesg-foundation/core/maintainability)
[![codecov](https://codecov.io/gh/mesg-foundation/core/branch/dev/graph/badge.svg)](https://codecov.io/gh/mesg-foundation/core)


MESG is a platform for the creation of efficient and easy-to-maintain applications that connect any and all technologies. 

Core is MESG's ultra-powerful communication and connection manager which manages the interaction of all connected services and applications so that applications can remain lightweight, yet feature packed.

To build an application, follow the Quick Start Guide, or if you'd like to help us build and maintain Core, refer to the Build From Source section below. 

# Contents

- [Quickstart](quickstart)
- [Service](service)
  - [Receiving Task](receiving-task)
  - [Submitting Event](submitting-event)
- [Architecture](architecture)
- [Examples](examples)
- [Roadmap](roadmap)
- [Community](community)
- [Contributing](contributing)


# Quick Start Guide

### 1. Download the CLI

First, download the CLI so you're able to interact with Core. You can either download the binaries directly from the [release page](https://github.com/mesg-foundation/core/releases/latest) then rename it to `mesg-core` and install it your path, or you can follow the installation process for your system in the [documentation](https://docs.mesg.tech/start-here/installing-the-cli)

### 2. Run MESG Core

MESG needs to have a daemon running to process all the different commands that you might need to execute. In order to start the daemon you can run:

```text
mesg-core start
```

### 3. Deploy a service

Next step is to deploy the service that your application will need. You can [create your own service](https://docs.mesg.tech/service/what-is-a-service), but for now, let's just use an existing one and deploy it.

```text
mesg-core deploy https://github.com/mesg-foundation/service-webhook
```

Let's deploy another one.

```text
mesg-core deploy https://github.com/mesg-foundation/service-invite-discord
```

Every time you deploy a service, the console will display the ID for the service you've just deployed.

### 4. Connect the services

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

### 5. Start the application

Start your application now like any node application:

```javascript
npm start
```

### 6. Test the application

Now we need to call the webhook in order to trigger the email, so let's do that with a curl command:

```text
curl -XPOST http://localhost:3000/webhook
```

You should now have an email in your inbox with your precious invitation to our Discord.

# Services

MESG depends heavily on services. These services are automatically built and ran inside Docker. You can connect anything you want, as long as it can run inside Docker \(as long as it can run on a computer\). If you need more details about how to connect dependencies to your service [check out the documentation](https://docs.mesg.tech/service/dockerize-the-service).

A service needs to implement two types of communications:receiving tasks and submitting events.

### Receiving Tasks

Tasks are designed to receive information both from Core and the Application that you run. Tasks can have multiple parameters as inputs and multiple outputs with varying data. You can visualize a task as a simple function that can return any kind of object.

You could have a task that receives a name as an input, and shows `success` as an output. This task factors the type of name with its probability like `{ "type": "female", "proabiliy": 92.34% }` but could also have an `error` output with a type of error like this: `{ "message": "This doesn't looks like a name" }`.

Check out the documentation for more information on how info how to create [tasks](https://docs.mesg.tech/service/listen-for-tasks).

### Submitting Events

Let's say you are working with a webserver. An event could be when there is a request with data in the payload, or it could be different events for each of the different routes of your API, or in a blockchain context, it could be when a smart contract emits an event.

For more info how to create your events, visit the [Emit an Event](https://docs.mesg.tech/service/emit-an-event) page.


# Architecture

\[\[ TODO: Add a nice graphic with the Application, the core and the services with the communication \]\]

# Examples

You can find a list of different examples and services that you can re-use [here](https://github.com/mesg-foundation/awesome)

# Roadmap

#### June 2018 - Core V1.0 Launched
Create your services and connect them together with your application through a single connection to Core, allowing Core to handle all communications and interoperability with any technology.

#### Q3 2018 - Rapid Deployment
No need to code your application anymore, just send a list of events with corresponding tasks within a simple configuration file to Core which will then execute tasks on your application’s behalf.

#### Q4 2018 - Beta Network
The decentralized beta Network means no coding or servers are necessary to run your applications. We will also launch an economy based on your participation in the Network.

#### Q3 2019 - Main Network
MESG launches its own blockchain Network providing for full scalability and a cheaper and faster user experience.

# Community

You can find us and other MESG users on Discord

[https://discordapp.com/invite/SaZ5HcE](https://discordapp.com/invite/SaZ5HcE)

Make sure to join and don't forget to introduce yourself and your project if you have one, also feel free to share good articles in the #neesfeed channel that might help other users and don't forget to stay nice and polite.

Don't forget also to check the [blog](https://medium.com/mesg), you might find interesting articles.

# Contributing

Contributions are more than welcome, in order to contribute please check at the [contribution guide](/blob/dev/CONTRIBUTING.md)

If you have any questions about it please reach out to us directly on [Discord](https://discordapp.com/invite/SaZ5HcE)
