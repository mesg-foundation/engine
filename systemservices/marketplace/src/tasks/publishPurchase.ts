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
    const txs = [...inputs.signedTransactions]
    // extract last transaction because we only want to get its logs
    const lastTransaction = txs.pop()
    for(const index in txs) {
      await web3.eth.sendSignedTransaction(txs[index])
    }
    const receipt = await web3.eth.sendSignedTransaction(lastTransaction)
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
