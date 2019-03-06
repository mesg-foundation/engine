import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { getServiceVersion } from "../contracts/version";
import { hexToAscii } from "../contracts/utils";

export default (
  contract: Marketplace,
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const sidHash = await contract.methods.hashToService(inputs.versionHash).call()
    if (hexToAscii(sidHash) === "") {
      throw new Error('version with hash ' + inputs.versionHash + ' does not exist')
    }
    const service = await contract.methods.services(sidHash).call()
    if (!await contract.methods.isServiceExist(service.sid).call()) {
      throw new Error('service with sid ' + hexToAscii(service.sid) + ' does not exist')
    }

    const authorizations = await Promise.all(inputs.addresses.map((address: string) => {
      return contract.methods.isAuthorized(service.sid, address).call()
    }))
    let authorized = false
    authorizations.forEach(authorization => {
      if (authorization === true) {
        authorized = true
      }
    })
    if (authorized === false) return outputs.success({
      authorized,
      sid: hexToAscii(service.sid),
    })

    const version = await getServiceVersion(contract, hexToAscii(service.sid), inputs.versionHash)
    if (version === undefined) {
      throw new Error('version with hash ' + inputs.versionHash + ' does not exist')
    }
    if (version.manifestData === undefined) {
      throw new Error('could not download manifest of version with hash ' + inputs.versionHash)
    }

    return outputs.success({
      authorized: authorized,
      sid: hexToAscii(service.sid),
      type: version.manifestData.service.deployment.type,
      source: version.manifestData.service.deployment.source,
    })
  }
  catch (error) {
    console.error('error in checkForDeployment', error)
    return outputs.error({ message: error.toString() })
  }
}
