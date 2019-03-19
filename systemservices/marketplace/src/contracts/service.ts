import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Service } from "../types/service";
import { getServiceVersions } from "./version";
import { getServiceOffers } from "./offer";
import { getServicePurchases } from "./purchase";
import { hexToAscii, parseTimestamp, asciiToHex } from "./utils";

const getAllServices = async (contract: Marketplace): Promise<Service[]> => {
  const servicesLength = new BigNumber(await contract.methods.servicesLength().call())
  const servicesPromise: Promise<Service|undefined>[] = []
  for (let i = new BigNumber(0); servicesLength.isGreaterThan(i); i = i.plus(1)) {
    servicesPromise.push(getServiceWithIndex(contract, i))
  }
  const services = await Promise.all(servicesPromise)
  return services.filter(service => service !== undefined) as Service[]
}

const getServiceWithIndex = async (contract: Marketplace, serviceIndex: BigNumber): Promise<Service|undefined> => {
  const sidHashed = await contract.methods.servicesList(serviceIndex.toString()).call()
  const service = await contract.methods.services(sidHashed).call()
  return getService(contract, hexToAscii(service.sid))
}

const getService = async (contract: Marketplace, sid: string): Promise<Service|undefined> => {
  const sidHex = asciiToHex(sid)
  if (!await contract.methods.isServiceExist(sidHex).call()) {
    return
  }
  const service = await contract.methods.service(sidHex).call()
  const [ versions, offers, purchases ] = await Promise.all([
    getServiceVersions(contract, sid),
    getServiceOffers(contract, sid),
    getServicePurchases(contract, sid),
  ])
  return {
    owner: service.owner,
    sid: sid,
    createTime: parseTimestamp(service.createTime),
    versions: versions,
    offers: offers,
    purchases: purchases,
  }
}

export { getAllServices, getService }
