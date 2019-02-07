import { service as MESG } from "mesg-js"
import Web3 from "web3"

import { newBlockEventEmitter } from "./newBlock"
import blockEvent from "./events/block"
import transactionAndLogEvent from "./events/transactionAndLog"

import listServices from "./tasks/listServices"
import createService from "./tasks/createService"
import sendSignedRawTx from "./tasks/sendSignedRawTx"

import _marketplaceABI from "./contracts/Marketplace.abi.json"
import { Marketplace } from "./contracts/Marketplace";
import { AbiItem } from "web3-utils/types";

const marketplaceABI = _marketplaceABI as AbiItem[]
const providerEndpoint = process.env.PROVIDER_ENDPOINT as string
const marketplaceAddress = process.env.MARKETPLACE_ADDRESS
const blockConfirmations = parseInt(<string>process.env.BLOCK_CONFIRMATIONS, 10)
const defaultGas = parseInt(<string>process.env.DEFAULT_GAS, 10)
const pollingTime = parseInt(<string>process.env.POLLING_TIME, 10)

let _defaultNonce = 0;
const defaultNonce = (): number => {
  _defaultNonce++
  return _defaultNonce
}

const main = async () => {
  const mesg = MESG()
  const web3 = new Web3(providerEndpoint)
  const marketplace = new web3.eth.Contract(marketplaceABI, marketplaceAddress) as Marketplace

  const defaultChainID = await web3.eth.net.getId()
  console.log('defaultChainID', defaultChainID)
  const defaultGasPrice = await web3.eth.getGasPrice()
  console.log('defaultGasPrice', defaultGasPrice)

  mesg.listenTask({
    listServices: listServices(marketplace),
    createService: createService(marketplace, defaultGas, defaultGasPrice.toString(), defaultChainID, defaultNonce),
    sendSignedRawTx: sendSignedRawTx(web3),
  })
  .on('error', error => console.error('catch listenTask', error))

  // const newBlock = await newBlockEventEmitter(web3, blockConfirmations, null, pollingTime)
  // newBlock.on('newBlock', blockNumber => {
  //   try {
  //     console.error('new block', blockNumber)

  //     blockEvent(mesg, web3, blockNumber)
  //     .catch(error => console.error('catch block event', error))
      
  //     transactionAndLogEvent(mesg, web3, blockNumber)
  //     .catch(error => console.error('catch transactionEvent', error))
  //   }
  //   catch (error) {
  //     console.error('catch newBlock on', error)
  //   }
  // })
}

try {
  main()
  .catch(error => console.error('catch promise', error))
} catch (error) {
  console.error('catch try', error)
}
