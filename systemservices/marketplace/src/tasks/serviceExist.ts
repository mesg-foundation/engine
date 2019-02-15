import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { getService } from "../contracts/service";
import { sha3 } from "../contracts/utils";

export default (contract: Marketplace) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const service = await getService(contract, sha3(inputs.sid))
    return outputs.success({ exist: service !== undefined })
  }
  catch (error) {
    console.error('error in serviceExist', error)
    return outputs.error({ message: error.toString() })
  }
}
