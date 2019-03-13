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
  const hash = hexToHash(await contract.methods.serviceHash(asciiToHex(sid), versionIndex.toString()).call())
  return getServiceVersion(contract, sid, hash)
}

const getServiceVersion = async (contract: Marketplace, sid: string, hash: string): Promise<Version|undefined> => {
  const sidHex = asciiToHex(sid)
  const hashHex = hashToHex(hash)
  if (!await contract.methods.isServiceVersionExist(sidHex, hashHex).call()) {
    return
  }
  const version = await contract.methods.serviceVersion(sidHex, hashHex).call()
  const manifestData = await getManifest(hexToAscii(version.manifestProtocol), hexToAscii(version.manifest))
  return {
    hash: hash,
    manifest: hexToAscii(version.manifest),
    manifestProtocol: hexToAscii(version.manifestProtocol),
    manifestData: manifestData,
    createTime: parseTimestamp(version.createTime),
  }
}

export { getServiceVersions, getServiceVersion }