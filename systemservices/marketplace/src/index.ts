import { service as MESG } from "mesg-js"
import Web3 from "web3"
const HttpProvider = require('web3-providers-http') // Wrong ts definition

import marketplaceABI from "./contracts/Marketplace.abi.json"
import { Marketplace } from "./contracts/Marketplace"
import ERC20ABI from "./contracts/ERC20.abi.json"
import { ERC20 } from "./contracts/ERC20"

import listenTasks from "./tasks"

const providerEndpoint = process.env.PROVIDER_ENDPOINT as string
const marketplaceAddress = process.env.MARKETPLACE_ADDRESS
const timeout = parseInt(<string>process.env.TIMEOUT, 10)
const ERC20Address = process.env.TOKEN_ADDRESS

const main = async () => {
  const mesg = MESG()
  const web3 = new Web3(new HttpProvider(providerEndpoint, { timeout }))
  const marketplace = new web3.eth.Contract(marketplaceABI, marketplaceAddress) as Marketplace
  const token = new web3.eth.Contract(ERC20ABI, ERC20Address) as ERC20

  const chainID = await web3.eth.net.getId()
  console.log('chainID', chainID)
  const defaultGasPrice = await web3.eth.getGasPrice()
  console.log('defaultGasPrice', defaultGasPrice)

  listenTasks(mesg, web3, marketplace, token, chainID, defaultGasPrice)

  console.log('service is ready and running')
}

try {
  main()
    .catch(error => console.error('catch promise', error))
} catch (error) {
  console.error('catch try', error)
}
