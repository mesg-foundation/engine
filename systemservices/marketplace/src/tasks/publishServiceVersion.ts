import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { asciiToHex, CreateTransaction, hashToHex } from "../contracts/utils";

export default (
  marketplace: Marketplace,
  createTransaction: CreateTransaction
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const transactionData = marketplace.methods.publishServiceVersion(
      asciiToHex(inputs.sid),
      asciiToHex(inputs.manifest),
      asciiToHex(inputs.manifestProtocol)
    ).encodeABI()
    return outputs.success(await createTransaction(marketplace, inputs, transactionData))
  }
  catch (error) {
    console.error('error in publishServiceVersion', error)
    return outputs.error({ message: error.toString() })
  }
}
