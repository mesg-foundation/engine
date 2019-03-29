import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import Web3 from "web3"
import { decodeLog } from "../contracts/utils";
import { eventHandlers } from "../events";
import { TransactionReceipt } from "web3/types";

export default (
  web3: Web3,
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const receipt: TransactionReceipt = inputs.signedTransactions.reduce(async (_: any, signedTx: string) => {
      return await web3.eth.sendSignedTransaction(signedTx)
    }, null)
    if (receipt.logs === undefined) throw new Error('receipt does not contain logs')
    const eventHandler = eventHandlers['ServicePurchased']
    const decodedLog = decodeLog(web3, eventHandler.abi, receipt.logs[receipt.logs.length - 1])
    const event = eventHandler.parse(decodedLog)
    return outputs.success(event)
  }
  catch (error) {
    console.error('error in publishPurchase', error)
    return outputs.error({ message: error.toString() })
  }
}
