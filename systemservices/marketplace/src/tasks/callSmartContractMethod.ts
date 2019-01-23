import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import Web3 from "web3"

export default (web3: Web3) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    // init web3 contract
    const contract: any = new web3.eth.Contract([inputs.methodAbi], inputs.contractAddress)

    // convert inputs parameters object to array
    const methodInputs = inputs.methodAbi.inputs.map((input: any) => {
      return inputs.inputs[input.name]
    })

    // generate the transaction's data
    let outputData = await contract.methods[inputs.methodAbi.name](...methodInputs).call()

    if (typeof outputData !== "object") {
      outputData = {
        0: outputData
      }
    }
    
    return await outputs.success({ outputs: outputData })
  }
  catch (error) {
    return await outputs.error({ message: error.toString() })
  }
}
