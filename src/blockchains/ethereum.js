const Web3 = require('web3')
const { InvalidBlockchainError } = require('../errors')
const { testConnection } = require('../utils')
const { emitRawBlock, emitRawTransaction } = require('../eventEmitter')

const type = 'ETHEREUM'
const endpoint = blockchain => process.env[`${type}_${blockchain.toUpperCase()}`]

module.exports = async ({ blockchain }) => {
  if (endpoint(blockchain) === undefined) {
    // We disable this blockchain throw the env variables but
    // if the env variable is present but empty this is an error
    return null
  }
  if (!endpoint(blockchain)) throw new InvalidBlockchainError(blockchain)

  const client = new Web3(new Web3.providers.HttpProvider(endpoint(blockchain)))

  await testConnection(() => client.isConnected(), `${type}/${blockchain}`)

  client.eth.filter('latest', (error, result) => {
    if (error) throw new Error('Error on watcher', error)
    const block = client.eth.getBlock(result, true)
    emitRawBlock({ type, blockchain, block })
    const receiptsBatch = client.createBatch()
    Promise.all(block.transactions
      .map((transaction, i) => new Promise((resolve, reject) => {
        receiptsBatch.add(client.eth.getTransactionReceipt.request(
          transaction.hash,
          (err, receipt) => err
            ? reject(err)
            : resolve({ ...transaction, ...receipt })
        ))
      })
    ))
      .then(transactions => transactions.forEach(transaction => {
        emitRawTransaction({ type, blockchain, block, transaction })
      }))
    receiptsBatch.execute()
  })
}
