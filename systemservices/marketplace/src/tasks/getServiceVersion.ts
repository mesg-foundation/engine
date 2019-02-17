import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { sha3 } from "../contracts/utils";
import { getServiceVersion } from "../contracts/version";

export default (contract: Marketplace) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const version = await getServiceVersion(contract, sha3(inputs.sid), inputs.hash)
    if (version === undefined) {
      throw new Error('version with hash ' + inputs.hash + ' and sid ' + inputs.sid + ' does not exist')
    }
    return outputs.success(version)
  }
  catch (error) {
    console.error('error in getServiceVersion', error)
    return outputs.error({ message: error.toString() })
  }
}
