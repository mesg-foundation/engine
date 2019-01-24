import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import Web3 from "web3"
import { Marketplace } from "../contracts/Marketplace"
import BigNumber from "bignumber.js"


const hexToAscii = (x: any) => Web3.utils.hexToAscii(x).replace(/\u0000/g, '')

interface Service {
  owner: string;
  sid: string;
  price: string;
  versions: Version[];
}

interface Version {
  hash: String;
  url: String;
}

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

const getServiceVersions = async (contract: Marketplace, serviceIndex: string): Promise<Version[]> => {
  const versionsCount = new BigNumber(await contract.methods.getServiceVersionsCount(serviceIndex).call())
  const versionsPromise: Promise<Version>[] = []
  for (let j = new BigNumber(0); versionsCount.isGreaterThan(j); j = j.plus(1)) {
    versionsPromise.push(getServiceVersion(contract, serviceIndex, j.toString()))
  }
  return await Promise.all(versionsPromise)
}

const getServiceVersion = async (contract: Marketplace, serviceIndex: string, versionIndex: string): Promise<Version> => {
  const version = await contract.methods.getServiceVersion(serviceIndex, versionIndex).call()
  return {
    hash: version.hash,
    url: hexToAscii(version.url),
  }
}

export default (contract: Marketplace) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const services = await getAllServices(contract)
    console.log('services', services)
    outputs.success({ services })
  }
  catch (error) {
    console.error('error im listServices', error)
    outputs.error({ message: error.toString() })
  }
}
