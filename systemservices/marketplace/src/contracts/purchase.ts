import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Purchase } from "../types/purchase";
import { parseTimestamp, asciiToHex } from "./utils";

const getServicePurchases = async (contract: Marketplace, sid: string): Promise<Purchase[]> => {
  const sidHex = asciiToHex(sid)
  if (!await contract.methods.isServiceExist(sidHex).call()) {
    return []
  }
  const purchasesLength = new BigNumber(await contract.methods.servicePurchasesLength(sidHex).call())
  const purchasesPromise: Promise<Purchase|undefined>[] = []
  for (let j = new BigNumber(0); purchasesLength.isGreaterThan(j); j = j.plus(1)) {
    purchasesPromise.push(getServicePurchaseWithIndex(contract, sid, j))
  }
  const purchases = await Promise.all(purchasesPromise)
  return purchases.filter(purchase => purchase !== undefined) as Purchase[]
}

const getServicePurchaseWithIndex = async (contract: Marketplace, sid: string, purchaseIndex: BigNumber): Promise<Purchase|undefined> => {
  const purchaser = await contract.methods.servicePurchaseAddress(asciiToHex(sid), purchaseIndex.toString()).call()
  return getServicePurchase(contract, sid, purchaser)
}

const getServicePurchase = async (contract: Marketplace, sid: string, purchaser: string): Promise<Purchase|undefined> => {
  const sidHex = asciiToHex(sid)
  if (!await contract.methods.isServicesPurchaseExist(sidHex, purchaser).call()) {
    return
  }
  const purchase = await contract.methods.servicePurchase(sidHex, purchaser).call()
  return {
    purchaser: purchaser,
    expire: parseTimestamp(purchase.expire),
    createTime: parseTimestamp(purchase.createTime),
  }
}

export { getServicePurchases, getServicePurchase }
