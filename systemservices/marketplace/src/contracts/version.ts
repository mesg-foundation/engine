import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Version } from "../types/version";
import { hexToString, parseTimestamp, stringToHex, hexToHash, hashToHex } from "./utils";

const getServiceVersions = async (contract: Marketplace, sid: string): Promise<Version[]> => {
  const sidHex = stringToHex(sid)
  if (!await contract.methods.isServiceExist(sidHex).call()) {
    throw new Error(`service ${sid} does not exist`)
  }
  const versionsLength = new BigNumber(await contract.methods.serviceVersionsLength(sidHex).call())
  const versionsPromise: Promise<Version>[] = []
  for (let j = new BigNumber(0); versionsLength.isGreaterThan(j); j = j.plus(1)) {
    versionsPromise.push(getServiceVersionWithIndex(contract, sid, j))
  }
  return Promise.all(versionsPromise)
}

const getServiceVersionWithIndex = async (contract: Marketplace, sid: string, versionIndex: BigNumber): Promise<Version> => {
  const versionHash = hexToHash(await contract.methods.serviceVersionHash(stringToHex(sid), versionIndex.toString()).call())
  return getServiceVersion(contract, versionHash)
}

const getServiceVersion = async (contract: Marketplace, versionHash: string): Promise<Version> => {
  const versionHashHex = hashToHex(versionHash)
  if (!await contract.methods.isServiceVersionExist(versionHashHex).call()) {
    throw new Error(`version ${versionHash} does not exist`)
  }
  const version = await contract.methods.serviceVersion(versionHashHex).call()
  return {
    versionHash: versionHash,
    manifest: hexToString(version.manifest),
    manifestProtocol: hexToString(version.manifestProtocol),
    createTime: parseTimestamp(version.createTime),
  }
}

export { getServiceVersions, getServiceVersion }
