import { service as MESG } from "mesg-js"
import Web3 from "web3"

import { newBlockEventEmitter } from "./newBlock"

import marketplaceABI from "./contracts/Marketplace.abi.json"
import { Marketplace } from "./contracts/Marketplace"
import ERC20ABI from "./contracts/ERC20.abi.json"
import { ERC20 } from "./contracts/ERC20"

import { eventHandlers } from "./events"

import preparePublishServiceVersion from "./tasks/preparePublishServiceVersion"
import publishPublishServiceVersion from "./tasks/publishPublishServiceVersion"
import getService from "./tasks/getService"
import prepareCreateServiceOffer from "./tasks/prepareCreateServiceOffer"
import listServices from "./tasks/listServices"
import purchase from "./tasks/purchase"
import sendSignedTransaction from "./tasks/sendSignedTransaction"
import publishCreateServiceOffer from "./tasks/publishCreateServiceOffer"
import isAuthorized from "./tasks/isAuthorized"
import { createTransactionTemplate } from "./contracts/utils";

const providerEndpoint = process.env.PROVIDER_ENDPOINT as string
const marketplaceAddress = process.env.MARKETPLACE_ADDRESS
const ERC20Address = process.env.TOKEN_ADDRESS
const blockConfirmations = parseInt(<string>process.env.BLOCK_CONFIRMATIONS, 10)
const defaultGas = parseInt(<string>process.env.DEFAULT_GAS, 10)
const pollingTime = parseInt(<string>process.env.POLLING_TIME, 10)

const main = async () => {
  const mesg = MESG()
  const web3 = new Web3(providerEndpoint)
  const marketplace = new web3.eth.Contract(marketplaceABI, marketplaceAddress) as Marketplace
  const token = new web3.eth.Contract(ERC20ABI, ERC20Address) as ERC20

  const chainID = await web3.eth.net.getId()
  console.log('chainID', chainID)
  const defaultGasPrice = await web3.eth.getGasPrice()
  console.log('defaultGasPrice', defaultGasPrice)

  const createTransaction = createTransactionTemplate(chainID, web3, defaultGas, defaultGasPrice)

  mesg.listenTask({
    listServices: listServices(marketplace),
    getService: getService(marketplace),
    preparePublishServiceVersion: preparePublishServiceVersion(marketplace, createTransaction),
    publishPublishServiceVersion: publishPublishServiceVersion(web3),
    prepareCreateServiceOffer: prepareCreateServiceOffer(marketplace, createTransaction),
    purchase: purchase(marketplace, token, createTransaction),
    sendSignedTransaction: sendSignedTransaction(web3),
    publishCreateServiceOffer: publishCreateServiceOffer(web3),
    isAuthorized: isAuthorized(marketplace),
  })
  .on('error', error => console.error('catch listenTask', error))

  const newBlock = await newBlockEventEmitter(web3, blockConfirmations, null, pollingTime)
  newBlock.on('newBlock', async blockNumber => {
    try {
      console.log('new block', blockNumber)
    
      const events = await marketplace.getPastEvents("allEvents", {
        fromBlock: blockNumber,
        toBlock: blockNumber,
      })
      events.forEach(async event => {
        // TODO: check if really async
        try {
          const eventHandler = eventHandlers[event.event]
          if (!eventHandler) {
            throw new Error('Event "'+event.event+'" is not implemented')
          }
          await mesg.emitEvent(eventHandler.mesgName, eventHandler.parse(event.returnValues))
        } catch(error) {
          return console.error('An error occurred during processing of event "'+event+ '". Error:', error)
        }
      })
    }
    catch (error) {
      console.error('catch newBlock on', error)
    }
  })

  console.log('service is ready and running')
}

try {
  main()
    .catch(error => console.error('catch promise', error))
} catch (error) {
  console.error('catch try', error)
}
