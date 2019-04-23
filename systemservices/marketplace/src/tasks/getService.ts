import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { getService } from "../contracts/service";
import { getManifest } from "../contracts/manifest";

export default (contract: Marketplace) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const service = await getService(contract, inputs.sid)
    const versionsWithManifest = await Promise.all(service.versions.map(async (version) => {
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
    return outputs.success({
      ...service,
      versions: versionsWithManifest,
    })
  }
  catch (error) {
    console.error('error in getService', error)
    return outputs.error({ message: error.message })
  }
}
