const Web3 = require('web3')
const { InvalidBlockchainError } = require('../errors')
const createClient = require('./createClient')

const onTransaction = client => callback => {
  client.eth
  .filter('latest', (error, result) => {
    if (error) throw new Error('Error on watcher', error)
    const block = client.eth.getBlock(result, true)
    block.transactions
      .forEach(transaction => {
        const receipt = client.eth.getTransactionReceipt(transaction.hash)
        callback(Object.assign({}, transaction, receipt), block)
      })
  })
}

module.exports = async chain => {
  const endpoint = process.env[`ETHEREUM_${chain.toUpperCase()}`]
  if (!endpoint) throw new InvalidBlockchainError(chain)

  const client = new Web3(new Web3.providers.HttpProvider(endpoint))
  
  return await createClient({
    type: 'ETHEREUM',
    network: chain.toUpperCase(),
    isConnected: () => client.isConnected(),
    onTransaction: onTransaction(client)
  })
}
