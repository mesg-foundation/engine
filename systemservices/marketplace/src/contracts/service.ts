import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Service } from "../types/service";
import { getServiceVersions } from "./version";

const getAllServices = async (contract: Marketplace): Promise<Service[]> => {
  const servicesCount = new BigNumber(await contract.methods.getServicesCount().call())
  const servicesPromise: Promise<Service>[] = []
  for (let i = new BigNumber(0); servicesCount.isGreaterThan(i); i = i.plus(1)) {
    servicesPromise.push(getService(contract, i.toString()))
  }
  return await Promise.all(servicesPromise)
}

const getService = async (contract: Marketplace, serviceIndex: string): Promise<Service> => {
  const [ service, versions ] = await Promise.all([
    contract.methods.services(serviceIndex).call(),
    getServiceVersions(contract, serviceIndex)
  ])
  return {
    owner: service.owner,
    sid: service.sid,
    price: service.price,
    versions: versions,
  }
}

export { getAllServices, getService }