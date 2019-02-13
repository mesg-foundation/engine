import { service as MESG } from "mesg-js"
import { TaskInputs, Service, EmitEventReply } from "mesg-js/lib/service"

import Web3 from "web3"
import { EventLog } from "web3/types"

import { newBlockEventEmitter } from "./newBlock"

import marketplaceABI from "./contracts/Marketplace.abi.json"
import { Marketplace } from "./contracts/Marketplace"

import createService from "./tasks/createService"
import createServiceOffer from "./tasks/createServiceOffer"
import createServiceVersion from "./tasks/createServiceVersion"
import disableServiceOffer from "./tasks/disableServiceOffer"
import listServices from "./tasks/listServices"
import purchase from "./tasks/purchase"
import sendSignedTransaction from "./tasks/sendSignedTransaction"
import transferServiceOwnership from "./tasks/transferServiceOwnership"
import isAuthorized from "./tasks/isAuthorized"
import serviceCreated = require("./events/serviceCreated");
import serviceOfferCreated = require("./events/serviceOfferCreated");
import serviceOfferDisabled = require("./events/serviceOfferDisabled");
import serviceOwnershipTransferred = require("./events/serviceOwnershipTransferred");
import serviceVersionCreated = require("./events/serviceVersionCreated");
import servicePurchased = require("./events/servicePurchased");

const providerEndpoint = process.env.PROVIDER_ENDPOINT as string
const marketplaceAddress = process.env.MARKETPLACE_ADDRESS
const blockConfirmations = parseInt(<string>process.env.BLOCK_CONFIRMATIONS, 10)
const defaultGas = parseInt(<string>process.env.DEFAULT_GAS, 10)
const pollingTime = parseInt(<string>process.env.POLLING_TIME, 10)

const eventHandlers: {[key: string]: (mesg: Service, event: EventLog) => Promise<EmitEventReply | Error>} = {
  'ServiceCreated': serviceCreated,
  'ServiceOfferCreated': serviceOfferCreated,
  'ServiceOfferDisabled': serviceOfferDisabled,
  'ServiceOwnershipTransferred': serviceOwnershipTransferred,
  'ServicePurchased': servicePurchased,
  'ServiceVersionCreated': serviceVersionCreated,
}

const main = async () => {
  const mesg = MESG()
  const web3 = new Web3(providerEndpoint)
  const contract = new web3.eth.Contract(marketplaceABI, marketplaceAddress) as Marketplace

  const chainID = await web3.eth.net.getId()
  console.log('chainID', chainID)
  const defaultGasPrice = await web3.eth.getGasPrice()
  console.log('defaultGasPrice', defaultGasPrice)

  const createTransaction = async (inputs: TaskInputs, data: string) => {
    return {
      chainID: chainID,
      nonce: await web3.eth.getTransactionCount(inputs.from),
      to: contract.options.address,
      gas: inputs.gas || defaultGas,
      gasPrice: inputs.gasPrice || defaultGasPrice,
      value: "0",
      data: data
    }
  }

  mesg.listenTask({
    listServices: listServices(contract),
    createService: createService(contract, createTransaction),
    createServiceOffer: createServiceOffer(contract, createTransaction),
    createServiceVersion: createServiceVersion(contract, createTransaction),
    disableServiceOffer: disableServiceOffer(contract, createTransaction),
    purchase: purchase(contract, createTransaction),
    transferServiceOwnership: transferServiceOwnership(contract, createTransaction),
    sendSignedTransaction: sendSignedTransaction(web3),
    isAuthorized: isAuthorized(contract),
  })
  .on('error', error => console.error('catch listenTask', error))

  // const newBlock = await newBlockEventEmitter(web3, blockConfirmations, null, pollingTime)
  // newBlock.on('newBlock', async blockNumber => {
  try {
    // console.error('new block', blockNumber)
    
    const events = await contract.getPastEvents("allEvents", {
      fromBlock: 5008717, // blockNumber,
      toBlock: 'latest', // blockNumber,
    })
    events.forEach(async event => {
      // TODO: check if really async
      try {
        if (!eventHandlers[event.event]) {
          return console.error('Event not implemented', event.event)
        }
        eventHandlers[event.event](mesg, event)
      } catch(error) {
        return console.error('An error occurred during processing of an event', event)
      }
    })
  }
  catch (error) {
    console.error('catch newBlock on', error)
  }
  // })

  console.log('service is ready and running')
}

try {
  main()
  .catch(error => console.error('catch promise', error))
} catch (error) {
  console.error('catch try', error)
}
