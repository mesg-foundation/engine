<p align="center">
  <img src="https://cdn.rawgit.com/mesg-foundation/core/dev/logo.svg" alt="MESG Core" height="120">
  <br/><br/>
</p>

[Website](https://mesg.com/) - [Docs](https://docs.mesg.com/) - [Chat](https://discordapp.com/invite/SaZ5HcE) - [Blog](https://medium.com/mesg)


[![CircleCI](https://img.shields.io/circleci/project/github/mesg-foundation/core.svg)](https://github.com/mesg-foundation/core)
[![Docker Pulls](https://img.shields.io/docker/pulls/mesg/core.svg)](https://hub.docker.com/r/mesg/core/)
[![Maintainability](https://api.codeclimate.com/v1/badges/86ad77f7c13cde40807e/maintainability)](https://codeclimate.com/github/mesg-foundation/core/maintainability)
[![codecov](https://codecov.io/gh/mesg-foundation/core/branch/dev/graph/badge.svg)](https://codecov.io/gh/mesg-foundation/core)


MESG is a platform for the creation of efficient and easy-to-maintain applications that connect any and all technologies. 

MESG Core is a communication and connection layer which manages the interaction of all connected services and applications so they can remain lightweight, yet feature packed.

To build an application, follow the [Quick Start Guide](#quick-start-guide)

If you'd like to build Services and share them with the community, go to the [Services](#services) section.

To help us build and maintain MESG Core, refer to the [Contribute](#contribute) section below.

# Contents

- [Quick Start Guide](#quick-start-guide)
- [Services](#services)
- [Architecture](#architecture)
- [Marketplace](#marketplace)
- [Roadmap](#roadmap)
- [Community](#community)
- [Contribute](#contribute)


# Quick Start Guide

This guide will show you steps-by-step how to create an application that sends a Discord invitation email when a webhook is called.

### 1. Installation

Run the following command in a console to install MESG Core:

```bash
bash <(curl -fsSL https://mesg.com/install)
```

You can also install it manually by following [this guide](https://docs.mesg.tech/start-here/installation#manual-installation).

### 2. Run MESG Core

MESG Core runs as a daemon. To start it, execute:

```bash
mesg-core start
```

### 3. Deploy the services

You need to deploy every service your application is using.

In this guide, the application is using 2 services.

Let's start by deploying the [webhook service](https://github.com/mesg-foundation/service-webhook):

```bash
mesg-core service deploy https://github.com/mesg-foundation/service-webhook
```

Now let's deploy the [invite discord service](https://github.com/mesg-foundation/service-discord-invitation):

```bash
mesg-core service deploy https://github.com/mesg-foundation/service-discord-invitation
```

Once the service is deployed, the console displays its Service ID. The Service ID is the unique way for the application to connect to the right service through MESG Core. You'll need to use them inside the application.

### 4. Create the application

Now that the services are up and running, let's create the application.

The application is using [NodeJS](https://nodejs.org) and [NPM](https://www.npmjs.com/).

Let's init the app and install the [MESG JS library](https://www.npmjs.com/package/mesg-js).

```bash
npm init && npm install --save mesg-js
```

Now, let's create an `index.js` file and with the following code:

```javascript
const MESG = require('mesg-js').application()

const webhook    = '__ID_SERVICE_WEBHOOK__' // To replace by the Service ID of the Webhook service
const invitation = '__ID_SERVICE_INVITATION_DISCORD__' // To replace by the Service ID of the Invite Discord service
const email      = '__YOUR_EMAIL_HERE__' // To replace by your email
const sendgridAPIKey = '__SENDGRID_API_KEY__' // To replace by your SendGrid API key. See https://app.sendgrid.com/settings/api_keys

MESG.whenEvent(
  { serviceID: webhook, filter: 'request' },
  { serviceID: invitation, taskKey: 'send', inputs: { email, sendgridAPIKey } },
)
```

Don't forget to replace the values `__ID_SERVICE_WEBHOOK__`, `__ID_SERVICE_INVITATION_DISCORD__`, `__YOUR_EMAIL_HERE__` and `__SENDGRID_API_KEY__`.

### 5. Start the application

Start your application like any node application:

```bash
node index.js
```

### 6. Test the application

Now let's give this super small application a try.

Let's trigger the webhook with the following command:

```bash
curl -XPOST http://localhost:3000/webhook
```

:tada: You should have received an email in your inbox with your precious invitation to our Discord.

# Services

Services are build and [shared by the community](https://github.com/mesg-foundation/awesome). They are small and reusable pieces of code that, when grouped together, allow developers to build incredible applications with ease.

You can develop a service for absolutely anything you want, as long as it can run inside Docker. Check the [documentation to create your own services](https://docs.mesg.com/service).

Services implement two types of communication: executing tasks and submitting events.

### Executing Tasks

Tasks have multiple input parameters and multiple outputs with varying data. A task is like a function with inputs and outputs.

Let's take an example of a task that sends a email:

The task accepts as inputs: `receiver`, `subject` and `body`.

The task could return 2 different outputs.

The first possible output is `success` with an empty object `{}` as data, meaning that the email has been sent with success

The second possible output is `error` with for eg, `{ "reason": "email invalid" }` as data.

This way, the application can easily check the type of output and react appropriately.

Check out the documentation for more information on [how to create tasks](https://docs.mesg.com/service/listen-for-tasks).

### Submitting Events

Services can also submit events to MESG Core. They allow two-way communication with MESG Core and Applications.

Let's say the service is a HTTP webserver. An event could be submitted when the webserver receives a request with the request's payload as the event's data. The service could also submit a specific event for every route of your API.

For more info how to create your events, visit the [Emit an Event](https://docs.mesg.com/service/emit-an-event) page.


# Architecture

![MESG Architecture](https://cdn.rawgit.com/mesg-foundation/core/dev/schema1.svg)

# Marketplace

We have a common place to post all community-developed Services and Applications. Check out the [curated list of Awesome Services and Applications](https://github.com/mesg-foundation/awesome) to participate.

Alternatively, you can also check out the [https://mesg.com/marketplace](Marketplace).

# Roadmap

#### June 2018 - Core V1.0 Launched
Create your services and connect them together with your application through a single connection to Core, allowing Core to handle all communications and interoperability with any technology. Services and applications can be shared with others in our Marketplace.

#### Q3 2018 - Rapid Deployment
No need to code your application anymore, just send a list of events with corresponding tasks within a simple configuration file to Core which will then execute tasks on your application’s behalf.

#### Q4 2018 - Beta Network
The decentralized beta Network means no need to run your applications and their services on your own computer, just deploy them on the Network. 

#### Q3 2019 - Main Network
MESG launches its own blockchain Network providing for full scalability and a cheaper and faster user experience.

# Community

You can find us and other MESG users on [Discord](https://discordapp.com/invite/SaZ5HcE).

Be sure to join, and don't forget to introduce yourself and your project if you have one.

Please feel free to share useful articles in the #newsfeed channel.

Also, be sure to check out the [blog](https://medium.com/mesg) to stay up-to-date with our articles.

# Contribute

Contributions are more than welcome. For more details on how to contribute, please check out the [contribution guide](/CONTRIBUTING.md).

If you have any questions, please reach out to us directly on [Discord](https://discordapp.com/invite/SaZ5HcE).
