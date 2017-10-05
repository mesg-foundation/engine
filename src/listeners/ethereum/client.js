const { InvalidBlockchainError } = require('../../errors')
const { testConnection } = require('../../utils')
const Web3 = require('web3')

const connectors = {}

const parseResult = data => ({
  blockId: data.blockNumber.toString(),
  fees: data.gasUsed.toString(),
  from: data.from,
  payload: data.args || {},
  to: data.to,
  transactionId: data.transactionHash,
  value: data.value
})

const getTransactionReceipt = (ethClient, transactionId) => new Promise((resolve, reject) => ethClient
  .getTransactionReceipt(transactionId, (err, res) => err ? reject(err) : resolve(res)))

const getTransaction = (ethClient, transactionId) => new Promise((resolve, reject) => ethClient
  .getTransaction(transactionId, (err, res) => err ? reject(err) : resolve(res)))

const merge = dataArray => dataArray.reduce((acc, value) => Object.assign(acc, value), {})

const handleEvent = web3Client => callback => (error, data) => error
  ? callback(error)
  : Promise.all([
    Promise.resolve(data),
    getTransactionReceipt(web3Client.eth, data.transactionHash),
    getTransaction(web3Client.eth, data.transactionHash)
  ])
    .then(merge)
    .then(parseResult)
    .then(data => callback(null, data))
    .catch(error => callback(error))

const nodeEndpoint = chain => {
  const endpoint = process.env[`ETHEREUM_${chain}`]
  if (!endpoint) throw new InvalidBlockchainError(chain)
  return endpoint
}

module.exports = async chain => {
  if (!connectors[chain]) {
    const endpoint = nodeEndpoint(chain)
    const web3Client = new Web3(new Web3.providers.HttpProvider(endpoint))
    await testConnection(() => web3Client.isConnected(), endpoint)
    web3Client.handleEvent = handleEvent(web3Client)
    connectors[chain] = web3Client
  }
  return connectors[chain]
}
