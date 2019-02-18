import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { getServiceVersion } from "../contracts/version";

export default (
  contract: Marketplace,
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const sidHash = await contract.methods.hashToService(inputs.hash).call()

    const authorizations = await Promise.all(inputs.addresses.map((address: string) => {
      return contract.methods.isAuthorized(sidHash, address).call()
    }))
    let authorized = false
    authorizations.forEach(authorization => {
      if (authorization === true) {
        authorized = true
      }
    })
    if (authorized === false) return outputs.success({ authorized })

    const version = await getServiceVersion(contract, sidHash, inputs.hash)
    if (version === undefined) {
      throw new Error('version with hash ' + inputs.hash + ' does not exist')
    }
    if (version.manifest === undefined) {
      throw new Error('could not download manifest of version with hash ' + inputs.hash)
    }

    return outputs.success({
      authorized: authorized,
      type: version.manifest.service.deployment.type,
      source: version.manifest.service.deployment.source,
    })
  }
  catch (error) {
    console.error('error in checkForDeployment', error)
    return outputs.error({ message: error.toString() })
  }
}
