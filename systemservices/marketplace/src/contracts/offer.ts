import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Offer } from "../types/service";
import { isValidNumber, fromUnit } from "./utils";

const getServiceOffers = async (contract: Marketplace, sidHash: string): Promise<Offer[]> => {
  const offersLength = new BigNumber(await contract.methods.servicesOffersLength(sidHash).call())
  if (!isValidNumber(offersLength)) {
    return []
  }
  const offersPromise: Promise<Offer|undefined>[] = []
  for (let j = new BigNumber(0); offersLength.isGreaterThan(j); j = j.plus(1)) {
    offersPromise.push(getServiceOffer(contract, sidHash, j))
  }
  const offers = await Promise.all(offersPromise)
  return offers.filter(offer => offer !== undefined) as Offer[]
}

const getServiceOffer = async (contract: Marketplace, sidHash: string, offerIndex: BigNumber): Promise<Offer|undefined> => {
  const offer = await contract.methods.servicesOffer(sidHash, offerIndex.toString()).call()
  const price = new BigNumber(offer.price)
  const duration = new BigNumber(offer.duration)
  if (!isValidNumber(price) || !isValidNumber(duration)) {
    return
  }
  return {
    index: offerIndex,
    price: fromUnit(price),
    duration: duration,
    active: offer.active,
  }
}

export { getServiceOffers, getServiceOffer }
