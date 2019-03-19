import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Offer } from "../types/offer";
import { fromUnit, parseTimestamp, asciiToHex } from "./utils";

const getServiceOffers = async (contract: Marketplace, sid: string): Promise<Offer[]> => {
  const sidHex = asciiToHex(sid)
  if (!await contract.methods.isServiceExist(sidHex).call()) {
    return []
  }
  const offersLength = new BigNumber(await contract.methods.serviceOffersLength(sidHex).call())
  const offersPromise: Promise<Offer|undefined>[] = []
  for (let j = new BigNumber(0); offersLength.isGreaterThan(j); j = j.plus(1)) {
    offersPromise.push(getServiceOffer(contract, sid, j))
  }
  const offers = await Promise.all(offersPromise)
  return offers.filter(offer => offer !== undefined) as Offer[]
}

const getServiceOffer = async (contract: Marketplace, sid: string, offerIndex: BigNumber): Promise<Offer|undefined> => {
  const sidHex = asciiToHex(sid)
  if (!await contract.methods.isServiceOfferExist(sidHex, offerIndex.toString()).call()) {
    return
  }
  const offer = await contract.methods.serviceOffer(sidHex, offerIndex.toString()).call()
  return {
    offerIndex: offerIndex,
    price: fromUnit(offer.price),
    duration: new BigNumber(offer.duration),
    active: offer.active,
    createTime: parseTimestamp(offer.createTime),
  }
}

export { getServiceOffers, getServiceOffer }
