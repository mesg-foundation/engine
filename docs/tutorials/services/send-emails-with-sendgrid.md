---
title: Send email with Sendgrid
description: 'Tutorial: how to use MESG to send an email through the Sendgrid API'
published_link: 'https://docs.mesg.com/tutorials/services/send-emails-with-sendgrid.html'
---

# Send emails with Sendgrid

## Introduction

Today we will learn how to create and use our first service.

We'll start with the example of an email provider where we will send emails through a MESG service.

You can access the final version of the [source code on GitHub](https://github.com/mesg-foundation/core/tree/master/docs/tutorials/services/send-email-with-sendgrid).

::: tip
If you haven't installed **MESG Core** yet, you can do so by running the command:

`bash <(curl -fsSL https://mesg.com/install)`
:::

MESG services are composed of two different parts:

* **Tasks:** actions that your service will execute
* **Events:** data that your service provides

For this tutorial, we will only focus on tasks, and will create a service that sends an email as a task.

To create a service with MESG, you need to create the base of the service. Run the command:

```bash
mesg-core service init
```

You should have a new folder with the name of your service that is composed of two different files:

* **Dockerfile:** A file that describes your Docker container and configures the environment for your service. We'll come back to this part in the [Dockerize your Service](send-emails-with-sendgrid.md#dockerize-your-service) section.
* **mesg.yml:** A file that contains all of the metadata of your Service. It gives some global descriptions but also includes the tasks and events that the Service can provide.

Let's add our first task.

## Send emails through MESG tasks

### Configure your task

Open the `mesg.yml` file and add a new task called `send` responsible for sending an email based on the following inputs:

* **from**: the account to send the email from
* **to**: the recipient of the email
* **subject**: the subject of the email
* **text**: the text of the email

Also this Service's task will return two different outputs:

* **success**: when the email is successfully sent
* **failure**: when an error occurs while trying to send the email \(this can be any kind of error, wrong email address, Sendgrid API down, etc.\)

To add this information into your Service, you can replace the `tasks: {}` in your `mesg.yml` file with the following:

```yaml
tasks: 
  send:
    inputs:
      from:
        type: String
      to:
        type: String
      subject:
        type: String
      text:
        type: String
    outputs:
      success:
        data:
          status:
            type: Number
      failure:
        data:
          message:
            type: String
```

::: warning
You might need to delete the **configuration: null** if it's present in your **mesg.yml** file
:::

### Code your task

The next step is to code your Service. To do this, you can use any language that you want. For this tutorial we will use **Javascript**.

Let's initialize our code for the Service. Run the following commands in the folder of your Service:

```bash
npm init -y                  # to initialize your project
npm install --save mesg-js   # to install the MESG library
```

When this is done, create and open a new file `index.js`

_Let's code !!!_

First we include and initialize the library to build a MESG Service.

```javascript
const MESG = require('mesg-js').service()
```

Secondly we'll need to create the different tasks defined in the `mesg.yml` as functions.

```javascript
const send = ({ from, to, subject, text }, { success, failure }) => { 
  // TODO later
}
```

For this task, the function will take all the inputs as a first parameter and all the outputs as a second parameter \(we are using [object destructuration](https://developer.mozilla.org/nl/docs/Web/JavaScript/Reference/Operatoren/Destructuring_assignment) here\).

The last step is to start listening for tasks from **MESG Core,** then react to those events.

```javascript
MESG.listenTask({ send })
```

Here we say that when there is the event `send` coming from MESG Core, we will execute the method `send` defined earlier.

The full file should look like this:

```javascript
const MESG = require('mesg-js').service()​

const send = ({ from, to, subject, text }, { success, failure }) => {
  // TODO later
}

​MESG.listenTask({ send })
```

This is a basic skeleton for a Service, but now we need to actually code the emails with Sendgrid. For this, we will use the Sendgrid library.

```bash
npm install --save @sendgrid/mail
```

Then require it in our file the same way we required it the `mesg-js` library.

```javascript
const sendgrid = require('@sendgrid/mail')
```

We had a `TODO` in our `send` task, let's code the business logic of this function using the Sendgrid library we just imported.

```javascript
const send = ({ from, to, subject, text }, { success, failure }) => {
  sendgrid.setApiKey('__CHANGE_WITH_YOUR_SENDGRID_API_KEY__')
  sendgrid.send({ from, to, subject, text })
    .then(([response, _]) => success({ status: response.statusCode }))
    .catch(e => failure({ message: e.toString() }))
}
```

This code is setting the API Key necessary to work with Sendgrid, then it sends an email with the parameters defined by our MESG Service. The result of this call is a promise that if successful, it will call the output `success` with the data `status` or if it fails, it will call the output `failure` with the `message` of the failure.

Now your final Service code should look like this:

```javascript
const MESG = require('mesg-js').service()
const sendgrid = require('@sendgrid/mail')

​const send = ({ from, to, subject, text }, { success, failure }) => {
  sendgrid.setApiKey('__CHANGE_WITH_YOUR_SENDGRID_API_KEY__')
  sendgrid.send({ from, to, subject, text })
    .then(([response, _]) => success({ status: response.statusCode }))
    .catch(e => failure({ message: e.toString() }))
}

​MESG.listenTask({ send })
```

::: warning
Don't forget to change the `__CHANGE_WITH_YOUR_SENDGRID_API_KEY__` with your own private Sendgrid API key that you can create here: [https://app.sendgrid.com/settings/api\_keys](https://app.sendgrid.com/settings/api_keys)​.
:::

Your Service is now ready for the second step.

## Dockerize your Service

This step is quite short and may not be necessary in the future. We need to edit the `Dockerfile` to make your Service compatible with Docker. In the case of a Javascript Service, the file will look like this:

```text
FROM node
WORKDIR /app
COPY ./package* ./
RUN npm install
COPY . .
CMD [ "node", "index.js" ]
```

With this file, your Service can now run in Docker, but don't worry, MESG will manage all this for you.

## Testing

The first step of testing is to make sure that the Service is valid by running:

```bash
mesg-core service validate
```

Your should have a message with `Service is valid`, if not, check the previous steps again; you probably missed something 🤔

Now that your Service is valid, let's create a test file to test your task. Create a `test.json` file is with all the inputs needed for your task.

```javascript
{
    "from": "sender@domain.tld",
    "to": "__YOUR_EMAIL__",
    "subject": "Test email",
    "text": "Hello world from MESG Service"
}
```

::: warning
Replace the **\_\_YOUR\_EMAIL\_\_** with your own email to test it. Don't worry, this is only done locally. We will not collect it 😀
:::

Now time to test it. Simply run the following command:

```bash
mesg-core service test --task send --data test.json
```

Don't worry, the first time you do this, it will take a bit of time because MESG Core is building your Service, but the subsequent times will be faster.

Once your Service starts, the `send` task will be executed and you should have the result in the console and your precious email in your mailbox.

## Usage

To be able to use your Service from an application, you will need to deploy it first. To do this, just run the command:

```bash
mesg-core service deploy
```

This will return an ID for the Service that you'll be able to use to start the Service, stop it, show the logs etc... or connect different events to the tasks of this Service, but that's for the next tutorial.

Get some rest now, you've done a good job creating your first Service with MESG. Of course you can always create more complicated Services. Make sure to check out the [documentation](https://docs.mesg.com) for more details.

## Final version of the source code

<repository url="https://github.com/mesg-foundation/core/tree/master/docs/tutorials/services/send-email-with-sendgrid"></repository>
