import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Offer } from "../types/offer";
import { fromUnit, parseTimestamp, asciiToHex } from "./utils";
import { requireServiceExist } from "./service";

const getServiceOffers = async (contract: Marketplace, sid: string): Promise<Offer[]> => {
  await requireServiceExist(contract, sid)
  const offersLength = new BigNumber(await contract.methods.serviceOffersLength(asciiToHex(sid)).call())
  const offersPromise: Promise<Offer>[] = []
  for (let j = new BigNumber(0); offersLength.isGreaterThan(j); j = j.plus(1)) {
    offersPromise.push(getServiceOffer(contract, sid, j))
  }
  return Promise.all(offersPromise)
}

const getServiceOffer = async (contract: Marketplace, sid: string, offerIndex: BigNumber): Promise<Offer> => {
  const sidHex = asciiToHex(sid)
  if (!await contract.methods.isServiceOfferExist(sidHex, offerIndex.toString()).call()) {
    throw new Error(`offer for service '${sid}' with offer index '${offerIndex.toString()}' does not exist`)
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
