import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { sha3 } from "../contracts/utils";

export default (
  contract: Marketplace,
  createTransaction: (inputs: TaskInputs, data: string) => Promise<any>
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const transactionData = contract.methods.createServiceOffer(
      sha3(inputs.sid),
      inputs.price,
      inputs.duration
    ).encodeABI()
    return outputs.success(await createTransaction(inputs, transactionData))
  }
  catch (error) {
    console.error('error in createServiceOffer', error)
    return outputs.error({ message: error.toString() })
  }
}
