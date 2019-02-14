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
  const sidHash = await contract.methods.servicesList(serviceIndex.toString()).call()
  const [ versions, offers, purchases ] = await Promise.all([
    getServiceVersions(contract, sidHash),
    getServiceOffers(contract, sidHash),
    getServicePurchases(contract, sidHash),
  ])
  const service = await contract.methods.services(sidHash).call()
  return {
    owner: service.owner,
    sid: hexToAscii(service.sid),
    sidHash: sidHash,
    versions: versions,
    offers: offers,
    purchases: purchases,
  }
}

export { getAllServices, getService }
