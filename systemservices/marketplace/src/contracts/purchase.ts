import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Purchase } from "../types/service";

const getServicePurchases = async (contract: Marketplace, sid: string): Promise<Purchase[]> => {
  const purchasesLength = new BigNumber((await contract.methods.servicesPurchasesListLength(sid).call()).length)
  const purchasesPromise: Promise<Purchase>[] = []
  for (let j = new BigNumber(0); purchasesLength.isGreaterThan(j); j = j.plus(1)) {
    purchasesPromise.push(getServicePurchase(contract, sid, j))
  }
  return await Promise.all(purchasesPromise)
}

const getServicePurchase = async (contract: Marketplace, sid: string, purchaseIndex: BigNumber): Promise<Purchase> => {
  const purchaser = (await contract.methods.servicesPurchasesList(sid, purchaseIndex.toString()).call()).purchaser
  const expire = (await contract.methods.servicesPurchases(sid, purchaser).call()).expire
  return {
    purchaser: purchaser,
    expire: new BigNumber(expire),
  }
}

export { getServicePurchases, getServicePurchase }