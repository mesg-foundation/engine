import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { asciiToHex } from "../contracts/utils";

export default (
  contract: Marketplace,
  createTransaction: (inputs: TaskInputs, data: string) => Promise<any>
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const transactionData = contract.methods.publishServiceVersion(
      asciiToHex(inputs.sid),
      inputs.versionHash,
      asciiToHex(inputs.manifest),
      asciiToHex(inputs.manifestProtocol)
    ).encodeABI()
    return outputs.success(await createTransaction(inputs, transactionData))
  }
  catch (error) {
    console.error('error in publishServiceVersion', error)
    return outputs.error({ message: error.toString() })
  }
}
