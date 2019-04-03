import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import Web3 from "web3"
import { extractEventFromLogs } from "../contracts/utils";
import { Marketplace } from "../contracts/Marketplace";
import { servicePurchased } from "../contracts/parseEvents";

export default (
  web3: Web3,
  marketplace: Marketplace,
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    if (inputs.signedTransactions.length === 2) {
      await web3.eth.sendSignedTransaction(inputs.signedTransactions[0])
    }
    const receipt = await web3.eth.sendSignedTransaction(inputs.signedTransactions[inputs.signedTransactions.length-1])
    if (receipt === null || receipt.logs === undefined) throw new Error('receipt does not contain logs')
    const decodedLog = extractEventFromLogs(web3, marketplace, 'ServicePurchased', receipt.logs)
    const event = servicePurchased(decodedLog)
    return outputs.success(event)
  }
  catch (error) {
    console.error('error in publishPurchase', error)
    return outputs.error({ message: error.toString() })
  }
}
