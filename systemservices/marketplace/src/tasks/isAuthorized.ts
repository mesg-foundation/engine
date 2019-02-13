import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"

export default (
  contract: Marketplace,
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const authorized = await contract.methods.isAuthorized(
      inputs.hashedSid,
      inputs.from
    ).call()
    return outputs.success({ authorized })
  }
  catch (error) {
    console.error('error in createService', error)
    return outputs.error({ message: error.toString() })
  }
}
