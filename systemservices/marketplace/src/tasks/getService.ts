import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { getService } from "../contracts/service";

export default (contract: Marketplace) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const service = await getService(contract, inputs.sid)
    return outputs.success(service)
  }
  catch (error) {
    console.error('error in getService', error)
    return outputs.error({ message: error.toString() })
  }
}
