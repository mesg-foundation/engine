import BigNumber from "bignumber.js"
import { Marketplace } from "./Marketplace"
import { Version } from "../types/service";
import { hexToAscii } from "./utils";
import { getManifest } from "./manifest";
import Contract from "web3/eth/contract";

const getServiceVersions = async (contract: Contract, sid: string): Promise<Version[]> => {
  const versionsLength = new BigNumber(await contract.methods.servicesVersionsListLength(sid).call())
  const versionsPromise: Promise<Version>[] = []
  for (let j = new BigNumber(0); versionsLength.isGreaterThan(j); j = j.plus(1)) {
    versionsPromise.push(getServiceVersion(contract, sid, j))
  }
  return await Promise.all(versionsPromise)
}

const getServiceVersion = async (contract: Contract, sid: string, versionIndex: BigNumber): Promise<Version> => {
  const versionHash = await contract.methods.servicesVersionsList(sid, versionIndex.toString()).call()
  const version = await contract.methods.servicesVersion(sid, versionHash).call()
  const manifest = await getManifest(hexToAscii(version.manifestProtocol), hexToAscii(version.manifest))
  return {
    hash: versionHash,
    manifestSource: hexToAscii(version.manifest),
    manifestProtocol: hexToAscii(version.manifestProtocol),
    manifest: manifest,
  }
}

export { getServiceVersions, getServiceVersion }