import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Offer } from "../types/offer";
import { fromUnit, parseTimestamp } from "./utils";

const getServiceOffers = async (contract: Marketplace, sidHash: string): Promise<Offer[]> => {
  if (!await contract.methods.isServiceExist(sidHash).call()) {
    return []
  }
  const offersLength = new BigNumber(await contract.methods.servicesOffersLength(sidHash).call())
  const offersPromise: Promise<Offer|undefined>[] = []
  for (let j = new BigNumber(0); offersLength.isGreaterThan(j); j = j.plus(1)) {
    offersPromise.push(getServiceOffer(contract, sidHash, j))
  }
  const offers = await Promise.all(offersPromise)
  return offers.filter(offer => offer !== undefined) as Offer[]
}

const getServiceOffer = async (contract: Marketplace, sidHash: string, offerIndex: BigNumber): Promise<Offer|undefined> => {
  if (!await contract.methods.isServiceOfferExist(sidHash, offerIndex.toString()).call()) {
    return
  }
  const offer = await contract.methods.servicesOffer(sidHash, offerIndex.toString()).call()
  return {
    index: offerIndex,
    price: fromUnit(offer.price),
    duration: new BigNumber(offer.duration),
    active: offer.active,
    createTime: parseTimestamp(offer.createTime),
  }
}

export { getServiceOffers, getServiceOffer }
