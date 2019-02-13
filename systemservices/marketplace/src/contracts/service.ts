import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Service } from "../types/service";
import { getServiceVersions } from "./version";
import { getServiceOffers } from "./offer";
import { getServicePurchases } from "./purchase";
import { hexToAscii } from "./utils";

const getAllServices = async (contract: Marketplace): Promise<Service[]> => {
  const servicesLength = new BigNumber(await contract.methods.servicesListLength().call())
  const servicesPromise: Promise<Service>[] = []
  for (let i = new BigNumber(0); servicesLength.isGreaterThan(i); i = i.plus(1)) {
    servicesPromise.push(getService(contract, i))
  }
  return await Promise.all(servicesPromise)
}

const getService = async (contract: Marketplace, serviceIndex: BigNumber): Promise<Service> => {
  const hashedSid = await contract.methods.servicesList(serviceIndex.toString()).call()
  const [ versions, offers, purchases ] = await Promise.all([
    getServiceVersions(contract, hashedSid),
    getServiceOffers(contract, hashedSid),
    getServicePurchases(contract, hashedSid),
  ])
  const service = await contract.methods.services(hashedSid).call()
  return {
    owner: service.owner,
    sid: hexToAscii(service.sid),
    hashedSid: hashedSid,
    versions: versions,
    offers: offers,
    purchases: purchases,
  }
}

export { getAllServices, getService }
