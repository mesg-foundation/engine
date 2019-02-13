import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Offer } from "../types/service";

const getServiceOffers = async (contract: Marketplace, hashedSid: string): Promise<Offer[]> => {
  const offersLength = new BigNumber(await contract.methods.servicesOffersLength(hashedSid).call())
  const offersPromise: Promise<Offer>[] = []
  for (let j = new BigNumber(0); offersLength.isGreaterThan(j); j = j.plus(1)) {
    offersPromise.push(getServiceOffer(contract, hashedSid, j))
  }
  return await Promise.all(offersPromise)
}

const getServiceOffer = async (contract: Marketplace, hashedSid: string, offerIndex: BigNumber): Promise<Offer> => {
  const offer = await contract.methods.servicesOffer(hashedSid, offerIndex.toString()).call()
  return {
    index: offerIndex,
    price: new BigNumber(offer.price),
    duration: new BigNumber(offer.duration),
    active: offer.active,
  }
}

export { getServiceOffers, getServiceOffer }
