import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { getService } from "../contracts/service";
import { sha3 } from "../contracts/utils";

export default (contract: Marketplace) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const service = await getService(contract, sha3(inputs.sid))
    if (service === undefined) {
      return outputs.error({
        message: 'service with sid ' + inputs.sid + ' does not exist',
        code: 'notFound',
      })
    }
    return outputs.success(service)
  }
  catch (error) {
    console.error('error in getService', error)
    return outputs.error({ message: error.toString() })
  }
}
