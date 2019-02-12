import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { asciiToHex } from "../contracts/utils"
import Contract from "web3/eth/contract";

export default (
  contract: Contract,
  createTransaction: (inputs: TaskInputs, data: string) => Promise<any>
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const transactionData = contract.methods.purchase(
      asciiToHex(inputs.sid),
      inputs.offerIndex
    ).encodeABI()
    return outputs.success(await createTransaction(inputs, transactionData))
  }
  catch (error) {
    console.error('error in purchase', error)
    return outputs.error({ message: error.toString() })
  }
}
