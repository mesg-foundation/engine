import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Service } from "../types/service";
import { getServiceVersions } from "./version";
import { getServiceOffers } from "./offer";
import { getServicePurchases } from "./purchase";
import { hexToAscii, parseTimestamp } from "./utils";

const getAllServices = async (contract: Marketplace): Promise<Service[]> => {
  const servicesLength = new BigNumber(await contract.methods.servicesListLength().call())
  const servicesPromise: Promise<Service|undefined>[] = []
  for (let i = new BigNumber(0); servicesLength.isGreaterThan(i); i = i.plus(1)) {
    servicesPromise.push(getServiceWithIndex(contract, i))
  }
  const services = await Promise.all(servicesPromise)
  return services.filter(service => service !== undefined) as Service[]
}

const getServiceWithIndex = async (contract: Marketplace, serviceIndex: BigNumber): Promise<Service|undefined> => {
  const sidHash = await contract.methods.servicesList(serviceIndex.toString()).call()
  return getService(contract, sidHash)
}

const getService = async (contract: Marketplace, sidHash: string): Promise<Service|undefined> => {
  const service = await contract.methods.services(sidHash).call()
  if (service.sid === null ||
      service.owner === "0x0000000000000000000000000000000000000000"
  ) {
    return
  }
  const [ versions, offers, purchases ] = await Promise.all([
    getServiceVersions(contract, sidHash),
    getServiceOffers(contract, sidHash),
    getServicePurchases(contract, sidHash),
  ])
  return {
    owner: service.owner.toLowerCase(),
    sid: hexToAscii(service.sid).toLowerCase(),
    sidHash: sidHash.toLowerCase(),
    versions: versions,
    offers: offers,
    purchases: purchases,
    createTime: parseTimestamp(service.createTime),
  }
}

export { getAllServices, getService }
