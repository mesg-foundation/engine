---
title: Listen to transfer of Ethereum ERC20 token
description: >-
  Tutorial: How to create a MESG Service that listens for the transfers of an
  Ethereum ERC20 token.
published_link: 'https://docs.mesg.com/tutorials/services/listen-for-transfers-of-an-ethereum-erc20-token.html'
---

# Listen for transfers of an Ethereum ERC20 token

## Introduction

In this tutorial, we will cover how to create a MESG Service that listens for the transfers of an Ethereum ERC20 token.

This Service will be developed with JavaScript and [Node.js](https://nodejs.org).  
We will use the library [Web3.js](https://web3js.readthedocs.io/en/1.0/) to interact with Ethereum through [Infura](https://infura.io/).

You can access the final version of the [source code on GitHub](https://github.com/mesg-foundation/core/tree/master/docs/tutorials/services/listen-to-transfer-of-ethereum-erc20-token).

::: tip
MESG Core should already be installed on your computer. If it isn't yet, [install it here](../../guide/start-here/installation.md).
:::

## Create the MESG service

It's time to create our MESG Service. First, open a terminal in your development folder and run the following command:

```bash
mesg-core service init
```

Then, answer the prompts with the following information:

```text
Name: service-ethereum-erc20
Description: Listen to transfers of an ERC20
```

The command should have created a `service-ethereum-erc20` folder containing `mesg.yml` and `Dockerfile` files.  
Leave these files intact; we'll return to them a bit later in this tutorial.

## Create a Node app

Let's initialize our node app.

First, make sure your terminal is pointed towards the newly-created `service-ethereum-erc20` folder.  
Then, run:

```bash
npm init -y
```

Let's install [Web3.js](https://web3js.readthedocs.io/en/1.0/):

```text
npm install --save web3
```

Then, create the `index.js` file in the Service folder.

### Initialize Web3 with Infura

The first step is to load Web3 and initialize it with Infura.  
Add the following code to the top of `index.js` :

```javascript
const Web3 = require('web3')
const web3 = new Web3('wss://mainnet.infura.io/_ws')
```

We are using the new WebSocket endpoint of Infura to listen to transfers. This endpoint is public, but is not production ready yet and it may change in the future. If you aren't able to listen for transfers at the end of this tutorial, please let us know.

### Specify the ERC20 contract

To listen to transfers of an ERC20, we'll have to direct both the contract ABI and its address to Web3.  
In this tutorial, we will use the TRON ERC20 token. You can find its ABI and address on [Etherscan](https://etherscan.io/address/0xf230b790e05390fc8295f4d3f60332c93bed42e2#code).  
For the simplicity of this tutorial, we will use only a small part of the ABI that exposes the transfers.

Create the file `erc20-abi.json` in the Service folder and copy/paste the following ABI:

<<< @/docs/tutorials/services/listen-to-transfer-of-ethereum-erc20-token/erc20-abi.json

Now, let's come back to `index.js` and initialize the contract with the ABI and the address. Add:

```javascript
const contract = new web3.eth.Contract(require('./erc20-abi.json'), "0xf230b790e05390fc8295f4d3f60332c93bed42e2")
```

### Listen for transfer events

We're finally ready to listen for transfers!

Web3, thanks to the ABI, gives us access to the contract neatly. Let's add the following code to `index.js` :

```javascript
contract.events.Transfer({fromBlock: 'latest'})
.on('data', event => {
  console.log('transfer', event)
})
```

Let's try it!

```bash
node index.js
```

::: warning
It might take a while to receive and display a transfer in the console. The events are received in real time, but if nobody is transferring this ERC20, you won't receive or see any events. You can go onto [Etherscan](https://etherscan.io/token/0xf230b790e05390fc8295f4d3f60332c93bed42e2) to see the transfers.
:::

Let's improve the output by showing only the useful information. Edit to match:

```javascript
contract.events.Transfer({fromBlock: 'latest'})
.on('data', event => {
  console.log('transfer', {
    blockNumber: event.blockNumber,
    transactionHash: event.transactionHash,
    from: event.returnValues.from,
    to: event.returnValues.to,
    value: event.returnValues.value / Math.pow(10, 6),
  })
})
```

::: tip
We have to `divide` value by `Math.pow(10, 6)` because of the number of decimals defined in this contract.
:::

Let's run it again:

```bash
node index.js
```

In the terminal, we should see something like:

```text
transfer { blockNumber: 5827612,
  transactionHash: '0x02019f4a80ad43019b8e69aed59e1dea0f03fb48d9df610686a1f590e8f6216d',
  from: '0x58993319Fc9e1b6cFAda8047B63a723Cceb1FfFE',
  to: '0x99f79B7A134db6e30d1b12F9Ee823339CaC0BA83',
  value: 11276.800815 }
transfer { blockNumber: 5827612,
  transactionHash: '0xf4a0aad5245417ae376cb9962c93bb9c599d8160cec49a5d82ba593033e657d2',
  from: '0x385dFF5650776188f4da150aA8b17a467812923b',
  to: '0xe8b69609342C337873cD20513e64be7FdE9feCf2',
  value: 100 }
```

::: tip Congratulation
You've built a Node app that listens in real-time to transfers of an ERC20 token!
:::

## Transform the node app to a MESG Service

Now, it's time to transform this node app to a fully-compatible MESG Service.

### Update mesg.yml

Let's add the event we want to emit to MESG Core to the `mesg.yml` file.

Now, clean the `mesg.yml` file, keeping only the keys: `name` and `description`. It should look like this:

```yaml
name: service-ethereum-erc20-tuto
description: Listen to transfers of an ERC20
```

Let's add the transfer event definition:

```yaml
events:
  transfer:
    data:
      blockNumber:
        type: Number
      transactionHash:
        type: String
      from:
        type: String
      to:
        type: String
      value:
        type: Number
```

This definition matches the JavaScript object we want to emit to MESG Core. You can refer to the [documentation](../../guide/service/service-file.md) for more information about the `mesg.yml` file.

### Add mesg-js lib

Let's install the `mesg-js` lib:

```bash
npm install --save mesg-js
```

Let's transform the file `index.js` to use the lib. Add at the top of the file:

```javascript
const MESG = require('mesg-js').service()
```

Replace `console.log` by `MESG.emitEvent`, like so:

```javascript
contract.events.Transfer({fromBlock: 'latest'})
.on('data', event => {
  MESG.emitEvent('transfer', {
    blockNumber: event.blockNumber,
    transactionHash: event.transactionHash,
    from: event.returnValues.from,
    to: event.returnValues.to,
    value: event.returnValues.value / Math.pow(10, 6),
  })
})
```

### Dockerize it

Let's update the `Dockerfile` to make our Service compatible with Docker. Because it is a Node.JS app, it's pretty simple:

<<< @/docs/tutorials/services/listen-to-transfer-of-ethereum-erc20-token/Dockerfile

Let's also create a `.dockerignore` file to ignore the `node_modules` from the build of the Service.

<<< @/docs/tutorials/services/listen-to-transfer-of-ethereum-erc20-token/.dockerignore

### Test the Service

It's time to test the Service with MESG!

Make sure MESG Core is running:

```bash
mesg-core start
```

Let's actually test the service! Run:

```bash
mesg-core service test
```

After the building and deployment steps, you should see that the Service has started:

```text
Service started
Listening for events from the service...
Listening for results from the service...
```

And finally, after a few seconds:

```text
2018/06/21 18:40:15 Receive event transfer : {"blockNumber":5828174,"from":"0x5B47bbA2F60AFb4870c3909a5b249F01E6d11BAe","to":"0x819B2368fa8781C4866237A0EA5E61Ec51492A32","transactionHash":"0x524751269a73294fa1fddf8fd584e40d51f4174df2a4ee8e081ea9a94ce7cc90","value":79}
```

::: tip Congratulation!
Hooray!!! 🎉 You finished building a MESG Service that listens for transfer of an ERC20 token!
:::

### Deploy the Service

To use this Service in your future application, you'll need to deploy it:

```text
mesg-core service deploy
```

This command returns the service's ID that will be required by your application.

## Final version of the source code

<repository url="https://github.com/mesg-foundation/core/tree/master/docs/tutorials/services/listen-to-transfer-of-ethereum-erc20-token"></repository>