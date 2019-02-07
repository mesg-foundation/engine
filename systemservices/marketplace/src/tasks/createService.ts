import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { asciiToHex } from "../contracts/utils"

export default (contract: Marketplace, defaultGas: number, defaultGasPrice: string, defaultChainID: number, defaultNonce: () => number) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const transactionData = contract.methods.createService(asciiToHex(inputs.sid), inputs.price).encodeABI()
    return outputs.success({
      chainID: inputs.chainID || defaultChainID,
      nonce: inputs.nonce || defaultNonce(),
      to: contract.options.address,
      gas: inputs.gas || defaultGas,
      gasPrice: inputs.gasPrice || defaultGasPrice,
      value: "0",
      data: transactionData
    })
  }
  catch (error) {
    console.error('error in createService', error)
    return outputs.error({ message: error.toString() })
  }
}
