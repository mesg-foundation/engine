import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import Web3 from "web3"
import { extractEventFromLogs } from "../contracts/utils";
import { serviceVersionCreated } from "../contracts/parseEvents";
import { Marketplace } from "../contracts/Marketplace";

export default (
  web3: Web3,
  marketplace: Marketplace,
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const receipt = await web3.eth.sendSignedTransaction(inputs.signedTransaction)
    if (receipt.logs === undefined) throw new Error('receipt does not contain logs')
    const decodedLog = extractEventFromLogs(web3, marketplace, 'ServiceVersionCreated', receipt.logs)
    const event = serviceVersionCreated(decodedLog)
    return outputs.success(event)
  }
  catch (error) {
    console.error('error in publishPublishServiceVersion', error)
    return outputs.error({ message: error.message })
  }
}
