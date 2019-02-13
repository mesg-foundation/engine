import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Version } from "../types/service";
import { hexToAscii } from "./utils";
import { getManifest } from "./manifest";

const getServiceVersions = async (contract: Marketplace, hashedSid: string): Promise<Version[]> => {
  const versionsLength = new BigNumber(await contract.methods.servicesVersionsListLength(hashedSid).call())
  const versionsPromise: Promise<Version>[] = []
  for (let j = new BigNumber(0); versionsLength.isGreaterThan(j); j = j.plus(1)) {
    versionsPromise.push(getServiceVersion(contract, hashedSid, j))
  }
  return await Promise.all(versionsPromise)
}

const getServiceVersion = async (contract: Marketplace, hashedSid: string, versionIndex: BigNumber): Promise<Version> => {
  const versionHash = await contract.methods.servicesVersionsList(hashedSid, versionIndex.toString()).call()
  const version = await contract.methods.servicesVersion(hashedSid, versionHash).call()
  const manifest = await getManifest(hexToAscii(version.manifestProtocol), hexToAscii(version.manifest))
  return {
    hash: versionHash,
    manifestSource: hexToAscii(version.manifest),
    manifestProtocol: hexToAscii(version.manifestProtocol),
    manifest: manifest,
  }
}

export { getServiceVersions, getServiceVersion }