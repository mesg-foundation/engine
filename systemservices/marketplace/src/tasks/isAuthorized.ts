import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { asciiToHex } from "../contracts/utils"
import Contract from "web3/eth/contract";

export default (
  contract: Contract
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const authorized = (await contract.methods.isAuthorized(
      asciiToHex(inputs.sid)
    ).call({ from: inputs.from })).authorized
    return outputs.success({ authorized })
  }
  catch (error) {
    console.error('error in createService', error)
    return outputs.error({ message: error.toString() })
  }
}
