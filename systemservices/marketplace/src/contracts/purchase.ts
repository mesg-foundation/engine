import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Purchase } from "../types/service";

const getServicePurchases = async (contract: Marketplace, hashedSid: string): Promise<Purchase[]> => {
  const purchasesLength = new BigNumber(await contract.methods.servicesPurchasesListLength(hashedSid).call())
  const purchasesPromise: Promise<Purchase>[] = []
  for (let j = new BigNumber(0); purchasesLength.isGreaterThan(j); j = j.plus(1)) {
    purchasesPromise.push(getServicePurchase(contract, hashedSid, j))
  }
  return await Promise.all(purchasesPromise)
}

const getServicePurchase = async (contract: Marketplace, hashedSid: string, purchaseIndex: BigNumber): Promise<Purchase> => {
  const purchaser = await contract.methods.servicesPurchasesList(hashedSid, purchaseIndex.toString()).call()
  const expire = await contract.methods.servicesPurchase(hashedSid, purchaser).call()
  return {
    purchaser: purchaser,
    expire: new BigNumber(expire),
  }
}

export { getServicePurchases, getServicePurchase }
