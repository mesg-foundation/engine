import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import Web3 from "web3"

const decodeLog = (web3: Web3) => (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    if (inputs.abi.anonymous === false) {
      // Remove first element because event is non-anonymous
      // https://web3js.readthedocs.io/en/1.0/web3-eth-abi.html#decodelog
      inputs.topics.splice(0, 1)
    }
    const decodedData = web3.eth.abi.decodeLog(inputs.abi.inputs, inputs.data, inputs.topics)
    return outputs.success({
      ...inputs,
      decodedData: decodedData,
    })
  }
  catch (error) {
    return outputs.error({ message: error.toString() })
  }
}
export default decodeLog