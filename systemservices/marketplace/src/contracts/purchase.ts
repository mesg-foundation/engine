import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Purchase } from "../types/purchase";
import { isValidNumber, parseTimestamp } from "./utils";
import purchase from "../tasks/purchase";

const getServicePurchases = async (contract: Marketplace, sidHash: string): Promise<Purchase[]> => {
  const purchasesLength = new BigNumber(await contract.methods.servicesPurchasesListLength(sidHash).call())
  if (!isValidNumber(purchasesLength)) {
    return []
  }
  const purchasesPromise: Promise<Purchase|undefined>[] = []
  for (let j = new BigNumber(0); purchasesLength.isGreaterThan(j); j = j.plus(1)) {
    purchasesPromise.push(getServicePurchase(contract, sidHash, j))
  }
  const purchases = await Promise.all(purchasesPromise)
  return purchases.filter(purchase => purchase !== undefined) as Purchase[]
}

const getServicePurchase = async (contract: Marketplace, sidHash: string, purchaseIndex: BigNumber): Promise<Purchase|undefined> => {
  const purchaser = (await contract.methods.servicesPurchasesList(sidHash, purchaseIndex.toString()).call()).toLowerCase()
  if (purchaser === '0x0000000000000000000000000000000000000000') {
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
