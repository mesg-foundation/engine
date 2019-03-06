import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { toUnit, asciiToHex } from "../contracts/utils";

export default (
  contract: Marketplace,
  createTransaction: (inputs: TaskInputs, data: string) => Promise<any>
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const transactionData = contract.methods.createServiceOffer(
      asciiToHex(inputs.sid),
      toUnit(inputs.price).toString(),
      inputs.duration
    ).encodeABI()
    return outputs.success(await createTransaction(inputs, transactionData))
  }
  catch (error) {
    console.error('error in createServiceOffer', error)
    return outputs.error({ message: error.toString() })
  }
}
