import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Purchase } from "../types/service";

const getServicePurchases = async (contract: Marketplace, sidHash: string): Promise<Purchase[]> => {
  const purchasesLength = new BigNumber(await contract.methods.servicesPurchasesListLength(sidHash).call())
  const purchasesPromise: Promise<Purchase>[] = []
  for (let j = new BigNumber(0); purchasesLength.isGreaterThan(j); j = j.plus(1)) {
    purchasesPromise.push(getServicePurchase(contract, sidHash, j))
  }
  return await Promise.all(purchasesPromise)
}

const getServicePurchase = async (contract: Marketplace, sidHash: string, purchaseIndex: BigNumber): Promise<Purchase> => {
  const purchaser = await contract.methods.servicesPurchasesList(sidHash, purchaseIndex.toString()).call()
  const expire = await contract.methods.servicesPurchase(sidHash, purchaser).call()
  return {
    purchaser: purchaser,
    expire: new BigNumber(expire),
  }
}

export { getServicePurchases, getServicePurchase }
