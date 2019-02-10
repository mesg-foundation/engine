import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Service } from "../types/service";
import { getServiceVersions } from "./version";

const getAllServices = async (contract: Marketplace): Promise<Service[]> => {
  const servicesLength = new BigNumber(await contract.methods.servicesListLength().call())
  const servicesPromise: Promise<Service>[] = []
  for (let i = new BigNumber(0); servicesLength.isGreaterThan(i); i = i.plus(1)) {
    servicesPromise.push(getService(contract, i.toString()))
  }
  return await Promise.all(servicesPromise)
}

const getService = async (contract: Marketplace, serviceIndex: string): Promise<Service> => {
  const [ serviceSid, versions ] = await Promise.all([
    contract.methods.servicesList(serviceIndex).call(),
    getServiceVersions(contract, serviceIndex)
  ])
  const serviceOwner = await contract.methods.services(serviceSid).call()
  return {
    owner: serviceOwner,
    sid: serviceSid,
    versions: versions
    // offers: offers,
    // purchasers: purchasers,
  }
}

export { getAllServices, getService }