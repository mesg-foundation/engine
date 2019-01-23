import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import Web3 from "web3"

export default (web3: Web3, defaultGasLimit: Number) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    // init web3 contract
    const contract: any = new web3.eth.Contract([inputs.methodAbi], inputs.contractAddress)

    // extract function abi from abi
    // const methodAbi = inputs.abi.filter((abi: any) => {
    //   return abi.type === 'function' && abi.name === inputs.method
    // })[0]
    // if (!methodAbi) {
    //   return outputs.error({ message: 'ABI does not contain method ' + inputs.method })
    // }

    // convert inputs parameters object to array
    const methodInputs = inputs.methodAbi.inputs.map((input: any) => {
      return inputs.inputs[input.name]
    })

    // generate the transaction's data
    const transactionData = contract.methods[inputs.methodAbi.name](...methodInputs).encodeABI()

    // send (and sign)
    const account = web3.eth.accounts.privateKeyToAccount(inputs.privateKey)
    const signedTransaction = await account.signTransaction({
      to: inputs.contractAddress,
      gas: inputs.gasLimit || defaultGasLimit, // optional
      gasPrice: inputs.gasPrice, // optional
      value: inputs.value || 0, // default 0 ETH
      data: transactionData
    })
    
    const transactionHash = await new Promise((resolve: (hash: String) => void, reject) => {
      return web3.eth.sendSignedTransaction(signedTransaction.rawTransaction)
        .on('transactionHash', resolve)
        .on('error', reject)
    })
    
    return await outputs.success({ transactionHash })
  }
  catch (error) {
    return await outputs.error({ message: error.toString() })
  }
}
