import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Version } from "../types/service";
import { hexToAscii } from "./utils";
import { getManifest } from "./manifest";

const getServiceVersions = async (contract: Marketplace, sidHash: string): Promise<Version[]> => {
  const versionsLength = new BigNumber(await contract.methods.servicesVersionsListLength(sidHash).call())
  const versionsPromise: Promise<Version>[] = []
  for (let j = new BigNumber(0); versionsLength.isGreaterThan(j); j = j.plus(1)) {
    versionsPromise.push(getServiceVersion(contract, sidHash, j))
  }
  return await Promise.all(versionsPromise)
}

const getServiceVersion = async (contract: Marketplace, sidHash: string, versionIndex: BigNumber): Promise<Version> => {
  const versionHash = await contract.methods.servicesVersionsList(sidHash, versionIndex.toString()).call()
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