import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Purchase } from "../types/purchase";
import { parseTimestamp, stringToHex } from "./utils";

const getServicePurchases = async (contract: Marketplace, sid: string): Promise<Purchase[]> => {
  const sidHex = stringToHex(sid)
  if (!await contract.methods.isServiceExist(sidHex).call()) {
    throw new Error(`service ${sid} does not exist`)
  }
  const purchasesLength = new BigNumber(await contract.methods.servicePurchasesLength(sidHex).call())
  const purchasesPromise: Promise<Purchase>[] = []
  for (let j = new BigNumber(0); purchasesLength.isGreaterThan(j); j = j.plus(1)) {
    purchasesPromise.push(getServicePurchaseWithIndex(contract, sid, j))
  }
  return Promise.all(purchasesPromise)
}

const getServicePurchaseWithIndex = async (contract: Marketplace, sid: string, purchaseIndex: BigNumber): Promise<Purchase> => {
  const purchaser = await contract.methods.servicePurchaseAddress(stringToHex(sid), purchaseIndex.toString()).call()
  return getServicePurchase(contract, sid, purchaser)
}

const getServicePurchase = async (contract: Marketplace, sid: string, purchaser: string): Promise<Purchase> => {
  const sidHex = stringToHex(sid)
  if (!await contract.methods.isServicesPurchaseExist(sidHex, purchaser).call()) {
    throw new Error(`purchase for service '${sid}' with purchase '${purchaser}' does not exist`)
  }
  const purchase = await contract.methods.servicePurchase(sidHex, purchaser).call()
  return {
    purchaser: purchaser,
    expire: parseTimestamp(purchase.expire),
    createTime: parseTimestamp(purchase.createTime),
  }
}

export { getServicePurchases, getServicePurchase }
