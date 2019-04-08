import Web3 from "web3"
import { Service } from "mesg-js/lib/service"

import { Marketplace } from "../contracts/Marketplace";
import { createTransactionTemplate } from "../contracts/utils";
import { ERC20 } from "../contracts/ERC20";

import preparePublishServiceVersion from "./preparePublishServiceVersion"
import publishPublishServiceVersion from "./publishPublishServiceVersion"
import getService from "./getService"
import prepareCreateServiceOffer from "./prepareCreateServiceOffer"
import listServices from "./listServices"
import preparePurchase from "./preparePurchase"
import publishPurchase from "./publishPurchase"
import publishCreateServiceOffer from "./publishCreateServiceOffer"
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
    preparePublishServiceVersion: preparePublishServiceVersion(marketplace, createTransaction),
    publishPublishServiceVersion: publishPublishServiceVersion(web3, marketplace),
    prepareCreateServiceOffer: prepareCreateServiceOffer(marketplace, createTransaction),
    preparePurchase: preparePurchase(marketplace, token, createTransaction),
    publishPurchase: publishPurchase(web3, marketplace),
    publishCreateServiceOffer: publishCreateServiceOffer(web3, marketplace),
    isAuthorized: isAuthorized(marketplace),
  })
  .on('error', error => console.error('catch listenTask', error))
}
