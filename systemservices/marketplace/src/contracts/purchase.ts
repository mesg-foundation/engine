import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Purchase } from "../types/purchase";
import { parseTimestamp } from "./utils";

const getServicePurchases = async (contract: Marketplace, sidHash: string): Promise<Purchase[]> => {
  if (!await contract.methods.isServiceExist(sidHash).call()) {
    return []
  }
  const purchasesLength = new BigNumber(await contract.methods.servicesPurchasesListLength(sidHash).call())
  const purchasesPromise: Promise<Purchase|undefined>[] = []
  for (let j = new BigNumber(0); purchasesLength.isGreaterThan(j); j = j.plus(1)) {
    purchasesPromise.push(getServicePurchaseWithIndex(contract, sidHash, j))
  }
  const purchases = await Promise.all(purchasesPromise)
  return purchases.filter(purchase => purchase !== undefined) as Purchase[]
}

const getServicePurchaseWithIndex = async (contract: Marketplace, sidHash: string, purchaseIndex: BigNumber): Promise<Purchase|undefined> => {
  const purchaser = (await contract.methods.servicesPurchasesList(sidHash, purchaseIndex.toString()).call()).toLowerCase()
  return getServicePurchase(contract, sidHash, purchaser)
}

const getServicePurchase = async (contract: Marketplace, sidHash: string, purchaser: string): Promise<Purchase|undefined> => {
  if (!await contract.methods.isServicesPurchaseExist(sidHash, purchaser).call()) {
    return
  }
  const purchase = await contract.methods.servicesPurchase(sidHash, purchaser).call()
  return {
    purchaser: purchaser,
    expire: parseTimestamp(purchase.expire),
    createTime: parseTimestamp(purchase.createTime),
  }
}

export { getServicePurchases, getServicePurchase }
