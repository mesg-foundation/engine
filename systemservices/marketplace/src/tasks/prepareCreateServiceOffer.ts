import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { toUnit, stringToHex, CreateTransaction } from "../contracts/utils";

export default (
  marketplace: Marketplace,
  createTransaction: CreateTransaction
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const transactionData = marketplace.methods.createServiceOffer(
      stringToHex(inputs.sid),
      toUnit(inputs.price),
      inputs.duration
    ).encodeABI()
    return outputs.success(await createTransaction(marketplace, inputs, transactionData))
  }
  catch (error) {
    console.error('error in prepareCreateServiceOffer', error)
    return outputs.error({ message: error.toString() })
  }
}
