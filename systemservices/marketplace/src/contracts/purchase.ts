import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Purchase } from "../types/service";
import Contract from "web3/eth/contract";

const getServicePurchases = async (contract: Contract, sid: string): Promise<Purchase[]> => {
  const purchasesLength = new BigNumber(await contract.methods.servicesPurchasesListLength(sid).call())
  const purchasesPromise: Promise<Purchase>[] = []
  for (let j = new BigNumber(0); purchasesLength.isGreaterThan(j); j = j.plus(1)) {
    purchasesPromise.push(getServicePurchase(contract, sid, j))
  }
  return await Promise.all(purchasesPromise)
}

const getServicePurchase = async (contract: Contract, sid: string, purchaseIndex: BigNumber): Promise<Purchase> => {
  const purchaser = await contract.methods.servicesPurchasesList(sid, purchaseIndex.toString()).call()
  const expire = await contract.methods.servicesPurchases(sid, purchaser).call()
  return {
    purchaser: purchaser,
    expire: new BigNumber(expire),
  }
}

export { getServicePurchases, getServicePurchase }
