import Web3 from "web3"
import { Service } from "mesg-js/lib/service"

import { Marketplace } from "../contracts/Marketplace";
import { createTransactionTemplate } from "../contracts/utils";
import { ERC20 } from "../contracts/ERC20";

import publishServiceVersion from "./publishServiceVersion"
import getService from "./getService"
import createServiceOffer from "./createServiceOffer"
import listServices from "./listServices"
import purchase from "./purchase"
import sendSignedTransaction from "./sendSignedTransaction"
import isAuthorized from "./isAuthorized"

export default async (
  mesg: Service,
  web3: Web3,
  marketplace: Marketplace,
  token: ERC20,
  chainID: number,
  defaultGas: number,
  defaultGasPrice: number
) => {
  const createTransaction = createTransactionTemplate(chainID, web3, defaultGas, defaultGasPrice)
  mesg.listenTask({
    listServices: listServices(marketplace),
    getService: getService(marketplace),
    publishServiceVersion: publishServiceVersion(marketplace, createTransaction),
    createServiceOffer: createServiceOffer(marketplace, createTransaction),
    purchase: purchase(marketplace, token, createTransaction),
    sendSignedTransaction: sendSignedTransaction(web3),
    isAuthorized: isAuthorized(marketplace),
  })
  .on('error', error => console.error('catch listenTask', error))
}
