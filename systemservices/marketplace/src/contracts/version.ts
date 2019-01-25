import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Version } from "../types/service";
import { hexToAscii } from "./utils";
import { getMetadata } from "./metadata";

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
  const metadata = await getMetadata(hexToAscii(version.url))
  return {
    hash: version.hash,
    metadata: metadata,
  }
}

export { getServiceVersions, getServiceVersion }