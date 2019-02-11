import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Service } from "../types/service";
import { getServiceVersions } from "./version";
import { getServiceOffers } from "./offer";
import { getServicePurchases } from "./purchase";
import { toBN, isBN } from "web3-utils";

const getAllServices = async (contract: Marketplace): Promise<Service[]> => {
  const servicesLength = new BigNumber((await contract.methods.servicesListLength().call()).length)
  const servicesPromise: Promise<Service>[] = []
  for (let i = new BigNumber(0); servicesLength.isGreaterThan(i); i = i.plus(1)) {
    servicesPromise.push(getService(contract, i))
  }
  return await Promise.all(servicesPromise)
}

const getService = async (contract: Marketplace, serviceIndex: BigNumber): Promise<Service> => {
  const sid = await contract.methods.servicesList(serviceIndex.toString()).call()
  const [ versions, offers, purchases ] = await Promise.all([
    getServiceVersions(contract, sid),
    getServiceOffers(contract, sid),
    getServicePurchases(contract, sid),
  ])
  const serviceOwner = (await contract.methods.services(sid).call()).owner
  return {
    owner: serviceOwner,
    sid: sid,
    versions: versions,
    offers: offers,
    purchases: purchases,
  }
}

export { getAllServices, getService }