# Quick Start Guide

This is a step-by-step guide to create an application that sends a Discord invitation email when a webhook is called.

## 1. Installation

Run the following command in a console to install MESG Core:

```bash
bash <(curl -fsSL https://mesg.com/install)
```

You can also install it manually by following [this guide](https://docs.mesg.tech/start-here/installation#manual-installation).

## 2. Run MESG Core

MESG Core runs as a daemon. To start it, execute:

```bash
mesg-core start
```

## 3. Deploy the services

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

## 4. Create the application

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

## 5. Start the application

Start your application like any node application:

```bash
node index.js
```

## 6. Test the application

Now let's give this super simple application a try.

Let's trigger the webhook with the following command:

```bash
curl -XPOST http://localhost:3000/webhook
```

:tada: You should have received an email in your inbox with your invitation to our Discord. Come join our community.

