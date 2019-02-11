import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Offer } from "../types/service";
import { hexToAscii } from "./utils";
import { getManifest } from "./manifest";

const getServiceOffers = async (contract: Marketplace, sid: string): Promise<Offer[]> => {
  const offersLength = new BigNumber((await contract.methods.servicesOffersLength(sid).call()).length)
  const offersPromise: Promise<Offer>[] = []
  for (let j = new BigNumber(0); offersLength.isGreaterThan(j); j = j.plus(1)) {
    offersPromise.push(getServiceOffer(contract, sid, j))
  }
  return await Promise.all(offersPromise)
}

const getServiceOffer = async (contract: Marketplace, sid: string, offerIndex: BigNumber): Promise<Offer> => {
  const offer = await contract.methods.servicesOffer(sid, offerIndex.toString()).call()
  return {
    index: new BigNumber(offerIndex),
    price: new BigNumber(offer.price),
    duration: new BigNumber(offer.duration),
    active: offer.active,
  }
}

export { getServiceOffers, getServiceOffer }