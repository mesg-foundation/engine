import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Version } from "../types/service";
import { hexToAscii, isValidNumber } from "./utils";
import { getManifest } from "./manifest";

const getServiceVersions = async (contract: Marketplace, sidHash: string): Promise<Version[]> => {
  const versionsLength = new BigNumber(await contract.methods.servicesVersionsListLength(sidHash).call())
  if (!isValidNumber(versionsLength)) {
    return []
  }
  const versionsPromise: Promise<Version|undefined>[] = []
  for (let j = new BigNumber(0); versionsLength.isGreaterThan(j); j = j.plus(1)) {
    versionsPromise.push(getServiceVersion(contract, sidHash, j))
  }
  const versions = await Promise.all(versionsPromise)
  return versions.filter(version => version !== undefined) as Version[]
}

const getServiceVersion = async (contract: Marketplace, sidHash: string, versionIndex: BigNumber): Promise<Version|undefined> => {
  const versionHash = await contract.methods.servicesVersionsList(sidHash, versionIndex.toString()).call()
  if (versionHash === '0x08c379a000000000000000000000000000000000000000000000000000000000') {
    return
  }
  const version = await contract.methods.servicesVersion(sidHash, versionHash).call()
  const manifest = await getManifest(hexToAscii(version.manifestProtocol), hexToAscii(version.manifest))
  return {
    hash: versionHash,
    manifestSource: hexToAscii(version.manifest),
    manifestProtocol: hexToAscii(version.manifestProtocol),
    manifest: manifest,
  }
}

export { getServiceVersions, getServiceVersion }