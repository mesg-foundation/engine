import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Offer } from "../types/service";

const getServiceOffers = async (contract: Marketplace, sidHash: string): Promise<Offer[]> => {
  const offersLength = new BigNumber(await contract.methods.servicesOffersLength(sidHash).call())
  const offersPromise: Promise<Offer>[] = []
  for (let j = new BigNumber(0); offersLength.isGreaterThan(j); j = j.plus(1)) {
    offersPromise.push(getServiceOffer(contract, sidHash, j))
  }
  return await Promise.all(offersPromise)
}

const getServiceOffer = async (contract: Marketplace, sidHash: string, offerIndex: BigNumber): Promise<Offer> => {
  const offer = await contract.methods.servicesOffer(sidHash, offerIndex.toString()).call()
  return {
    index: offerIndex,
    price: new BigNumber(offer.price),
    duration: new BigNumber(offer.duration),
    active: offer.active,
  }
}

export { getServiceOffers, getServiceOffer }
