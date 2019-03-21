import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import Web3 from "web3"

export default (web3: Web3) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const receipt = await web3.eth.sendSignedTransaction(inputs.signedTransaction)
    return outputs.success(receipt)
  }
  catch (error) {
    console.error('error in sendSignedRawTx', error)
    return outputs.error({ message: error.toString() })
  }
}
