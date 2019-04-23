import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Version } from "../types/version";
import { hexToString, parseTimestamp, stringToHex, hexToHash, hashToHex, keccak256 } from "./utils";
import { requireServiceExist } from "./service";
import * as assert from "assert";

const getServiceVersions = async (contract: Marketplace, sid: string): Promise<Version[]> => {
  await requireServiceExist(contract, sid)
  const versionsLength = new BigNumber(await contract.methods.serviceVersionsLength(stringToHex(sid)).call())
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
  assert.ok(await isVersionExist(contract, versionHash), `version does not exist`)
  const version = await contract.methods.serviceVersion(hashToHex(versionHash)).call()
  return {
    versionHash: versionHash,
    manifest: hexToString(version.manifest),
    manifestProtocol: hexToString(version.manifestProtocol),
    createTime: parseTimestamp(version.createTime),
  }
}

const isVersionExist = async (contract: Marketplace, versionHash: string): Promise<boolean> => {
  return contract.methods.isServiceVersionExist(hashToHex(versionHash)).call()
}

const computeVersionHash = (from: string, sid: string, manifest: string, manifestProtocol: string) => {
  return hexToHash(keccak256(from, stringToHex(sid), stringToHex(manifest), stringToHex(manifestProtocol)))
}

export {
  getServiceVersions,
  getServiceVersion,
  isVersionExist,
  computeVersionHash,
}
