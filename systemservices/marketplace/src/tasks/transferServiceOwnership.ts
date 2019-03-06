import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { asciiToHex, CreateTransaction } from "../contracts/utils";

export default (
  marketplace: Marketplace,
  createTransaction: CreateTransaction
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const transactionData = marketplace.methods.transferServiceOwnership(
      asciiToHex(inputs.sid),
      inputs.newOwner
    ).encodeABI()
    return outputs.success(await createTransaction(marketplace, inputs, transactionData))
  }
  catch (error) {
    console.error('error in transferServiceOwnership', error)
    return outputs.error({ message: error.toString() })
  }
}
