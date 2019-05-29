import { TaskInputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { getService } from "../contracts/service";
import { getManifest } from "../contracts/manifest";
import { getServiceVersions } from "../contracts/version";
import { getServiceOffers } from "../contracts/offer";
import { getServicePurchases } from "../contracts/purchase";

export default (contract: Marketplace) => async (inputs: TaskInputs): Promise<object> => {
  const [service, versions, offers, purchases] = await Promise.all([
    getService(contract, inputs.sid),
    getServiceVersionWithManifest(contract, inputs.sid),
    getServiceOffers(contract, inputs.sid),
    getServicePurchases(contract, inputs.sid),
  ])
  return {
    ...service,
    offers,
    purchases,
    versions,
  }
}

const getServiceVersionWithManifest = async (contract: Marketplace, sid: string) => {
  const versions = await getServiceVersions(contract, sid)
  return await Promise.all(versions.map(async (version) => {
    let manifestData = null
    try {
      manifestData = await getManifest(version.manifestProtocol, version.manifest)
    }
    catch (error) {
      console.warn('error getManifest', error.toString())
    }
    return {
      ...version,
      manifestData
    }
  }))
}
