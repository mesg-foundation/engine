import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import Web3 from "web3"
import { decodeLog } from "../contracts/utils";
import { eventHandlers } from "../events";

export default (
  web3: Web3,
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    if (inputs.logs.length === 0) throw new Error('wrong number of logs')
    const eventHandler = eventHandlers['ServiceOfferCreated']
    const decodedLog = decodeLog(web3, eventHandler.abi, inputs.logs[inputs.logs.length - 1])
    const event = eventHandler.parse(decodedLog)
    return outputs.success(event)
  }
  catch (error) {
    console.error('error in decodeLogCreateServiceOffer', error)
    return outputs.error({ message: error.toString() })
  }
}
