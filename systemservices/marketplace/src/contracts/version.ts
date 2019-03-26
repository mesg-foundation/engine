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
  const versions: Version[] = []
  for (let j = new BigNumber(0); versionsLength.isGreaterThan(j); j = j.plus(1)) {
    try {
      const version = normalizeVersion(await getServiceVersionWithIndex(contract, sid, j))
      if (validVersion(version)) {
        versions.push(version)
      }
    } catch(e) {
      console.warn(e.toString())
    }
  }
  return versions
}

const getServiceVersionWithIndex = async (contract: Marketplace, sid: string, versionIndex: BigNumber): Promise<Version> => {
  const versionHash = hexToHash(await contract.methods.serviceVersionHash(asciiToHex(sid), versionIndex.toString()).call())
  return getServiceVersion(contract, versionHash)
}

const getServiceVersion = async (contract: Marketplace, versionHash: string): Promise<Version> => {
  const versionHashHex = hashToHex(versionHash)
  if (!await contract.methods.isServiceVersionExist(versionHashHex).call()) {
    throw new Error('Service does not exists')
  }
  const version = await contract.methods.serviceVersion(versionHashHex).call()
  const manifest = await getManifest(hexToAscii(version.manifestProtocol), hexToAscii(version.manifest))
  if (!manifest) {
    throw new Error('Unable to get the manifest')
  }
  return {
    versionHash: versionHash,
    manifest: manifest,
    createTime: parseTimestamp(version.createTime),
  }
}

const validVersion = (version: Version): boolean => {
  return version.versionHash != '' &&
    !!version.manifest &&
    !!version.manifest.service &&
    !!version.manifest.service.hash &&
    !!version.manifest.service.hashVersion &&
    !!version.manifest.service.definition &&
    !!version.manifest.service.deployment
}

const normalizeVersion = (version: any): Version => {
  return {
    createTime: version.createTime,
    versionHash: version.versionHash,
    manifest: {
      version: version.manifest.version.toString(),
      service: {
        definition: version.manifest.service.definition,
        deployment: {
          type: version.manifest.service.deployment.type || 'ipfs',
          source: version.manifest.service.deployment.source,
        },
        hash: version.manifest.service.hash,
        hashVersion: version.manifest.service.hashVersion,
        readme: version.manifest.service.readme || '',
      }
    }
  }
}

export { getServiceVersions, getServiceVersion }
