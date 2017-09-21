const Web3 = require('web3')

const connectors = {}

const promisifyFunctions = (object, functions) => {
  
  object.getTransactionReceiptPromise = transactionId => new Promise((resolve, reject) => object
    .getTransactionReceipt(transactionId, (err, res) => err ? reject(err) : resolve(res)))

  object.getTransactionPromise = transactionId => new Promise((resolve, reject) => object
    .getTransaction(transactionId, (err, res) => err ? reject(err) : resolve(res)))
  return object
}

module.exports = chain => {
  if (!connectors[chain]) {
    const web3Client = new Web3(new Web3.providers.HttpProvider(process.env.NODE_ADDRESS))
    connectors[chain] = promisifyFunctions(web3Client.eth)
  }
  return connectors[chain]
}