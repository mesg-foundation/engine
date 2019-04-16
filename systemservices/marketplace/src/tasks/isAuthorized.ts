import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { getServiceVersion } from "../contracts/version";
import { hexToString, stringToHex, hashToHex, hexToHash } from "../contracts/utils";
import BigNumber from "bignumber.js";
import { requireServiceExist } from "../contracts/service";

export default (
  contract: Marketplace,
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    let sid = ""
    let versionHash = ""
    if (inputs.versionHash) {
      versionHash = inputs.versionHash
      // get service from version hash
      const sidHash = await contract.methods.versionHashToService(hashToHex(versionHash)).call()
      if (hexToString(sidHash) === "") {
        throw new Error('service with hash ' + versionHash + ' does not exist')
      }
      const service = await contract.methods.services(sidHash).call()
      sid = hexToString(service.sid)
    }
    else if (inputs.sid) {
      // get version hash from sid
      sid = inputs.sid
      const versionsLength = new BigNumber(await contract.methods.serviceVersionsLength(stringToHex(sid)).call())
      if (versionsLength.isEqualTo(0)) {
        throw new Error('service with sid ' + sid + ' does not have any version')
      }
      versionHash = hexToHash(await contract.methods.serviceVersionHash(stringToHex(sid), versionsLength.minus(1).toString()).call())
    }
    else {
      throw new Error('input should have sid or hash set')
    }

    await requireServiceExist(contract, sid)

    // check if at least one of the provided addresses is authorized
    const authorizations = await Promise.all(inputs.addresses.map((address: string) => {
      return contract.methods.isAuthorized(stringToHex(sid), address).call()
    }) as Promise<boolean>[])
    const authorized = authorizations.reduce((p, c) => p || c, false)
    if (!authorized) {
      return outputs.success({
        authorized,
        sid: sid,
      })
    }

    // get version's manifest data
    const version = await getServiceVersion(contract, versionHash)
    if (version.manifestData === null) {
      throw new Error('could not download manifest of version with hash ' + versionHash)
    }

    return outputs.success({
      authorized: authorized,
      sid: sid,
      type: version.manifestData.service.deployment.type,
      source: version.manifestData.service.deployment.source,
    })
  }
  catch (error) {
    console.error('error in isAuthorized', error)
    return outputs.error({ message: error.toString() })
  }
}
