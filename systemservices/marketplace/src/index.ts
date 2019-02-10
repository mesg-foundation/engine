import { service as MESG } from "mesg-js"
import Web3 from "web3"

import { newBlockEventEmitter } from "./newBlock"
import blockEvent from "./events/block"
import transactionAndLogEvent from "./events/transactionAndLog"

import listServices from "./tasks/listServices"
import createService from "./tasks/createService"
import sendSignedTransaction from "./tasks/sendSignedTransaction"

import _marketplaceABI from "./contracts/Marketplace.abi.json"
import { Marketplace } from "./contracts/Marketplace";
import { AbiItem } from "web3-utils/types";
import { TaskInputs } from "mesg-js/lib/service";

const marketplaceABI = _marketplaceABI as AbiItem[]
const providerEndpoint = process.env.PROVIDER_ENDPOINT as string
const marketplaceAddress = process.env.MARKETPLACE_ADDRESS
const blockConfirmations = parseInt(<string>process.env.BLOCK_CONFIRMATIONS, 10)
const defaultGas = parseInt(<string>process.env.DEFAULT_GAS, 10)
const pollingTime = parseInt(<string>process.env.POLLING_TIME, 10)

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
      nonce: await web3.eth.getTransactionCount(inputs.address),
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
    sendSignedTransaction: sendSignedTransaction(web3),
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
