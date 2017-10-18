const Web3 = require('web3')
const Logger = require('../logger')
const { InvalidBlockchainError } = require('../errors')
const createClient = require('./createClient')

const onTransaction = (client, chain) => callback => client.eth
  .filter('latest', (error, result) => {
    if (error) throw new Error('Error on watcher', error)
    const block = client.eth.getBlock(result, true)
    Logger.debug(`==> [ETHEREUM ${chain}] BLOCK ${block.number} (${block.transactions.length} tx)`)
    const receiptsBatch = client.createBatch()
    Promise.all(block.transactions
      .map((transaction, i) => new Promise((resolve, reject) => {
        receiptsBatch.add(client.eth.getTransactionReceipt.request(
          transaction.hash,
          (err, receipt) => err ? reject(err) : resolve(Object.assign({}, transaction, receipt))
        ))
      })
    ))
      .then(transactions => transactions.forEach(transaction => {
        callback(transaction, block)
      }))
    receiptsBatch.execute()
  })

module.exports = async chain => {
  const endpoint = process.env[`ETHEREUM_${chain.toUpperCase()}`]
  if (!endpoint) throw new InvalidBlockchainError(chain)

  const client = new Web3(new Web3.providers.HttpProvider(endpoint))

  return createClient({
    type: 'ETHEREUM',
    network: chain.toUpperCase(),
    isConnected: () => client.isConnected(),
    onTransaction: onTransaction(client, chain)
  })
}
