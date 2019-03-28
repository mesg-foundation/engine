import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Version } from "../types/version";
import { hexToAscii, parseTimestamp, asciiToHex, hexToHash, hashToHex } from "./utils";
import { getManifest } from "./manifest";

const getServiceVersions = async (contract: Marketplace, sid: string): Promise<Version[]> => {
  const sidHex = asciiToHex(sid)
  if (!await contract.methods.isServiceExist(sidHex).call()) {
    return []
  }
  const versionsLength = new BigNumber(await contract.methods.serviceVersionsLength(sidHex).call())
  const versionsPromise: Promise<Version|undefined>[] = []
  for (let j = new BigNumber(0); versionsLength.isGreaterThan(j); j = j.plus(1)) {
    versionsPromise.push(getServiceVersionWithIndex(contract, sid, j))
  }
  const versions = await Promise.all(versionsPromise)
  return versions.filter(version => version !== undefined) as Version[]
}

const getServiceVersionWithIndex = async (contract: Marketplace, sid: string, versionIndex: BigNumber): Promise<Version|undefined> => {
  const versionHash = hexToHash(await contract.methods.serviceVersionHash(asciiToHex(sid), versionIndex.toString()).call())
  return getServiceVersion(contract, versionHash)
}

const getServiceVersion = async (contract: Marketplace, versionHash: string): Promise<Version|undefined> => {
  const versionHashHex = hashToHex(versionHash)
  if (!await contract.methods.isServiceVersionExist(versionHashHex).call()) {
    return
  }
  const version = await contract.methods.serviceVersion(versionHashHex).call()
  let manifestData = null
  try {
    manifestData = await getManifest(hexToAscii(version.manifestProtocol), hexToAscii(version.manifest))
  }
  catch (error) {
    console.warn(error.toString())
  }
  return {
    versionHash: versionHash,
    manifest: hexToAscii(version.manifest),
    manifestProtocol: hexToAscii(version.manifestProtocol),
    manifestData: manifestData,
    createTime: parseTimestamp(version.createTime),
  }
}

export { getServiceVersions, getServiceVersion }
