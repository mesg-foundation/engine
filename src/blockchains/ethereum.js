const Web3 = require('web3')
const Logger = require('../logger')
const { InvalidBlockchainError } = require('../errors')
const { emitRawBlock, emitRawTransaction } = require('../eventEmitter')

const type = 'ETHEREUM'
const endpoint = blockchain => process.env[`${type}_${blockchain.toUpperCase()}`]

const handlerNewBlock = (client, blockchain) => async result => {
  try {
    const block = await client.eth.getBlock(result.number, true)
    emitRawBlock({ type, blockchain, block })
    if (block.transactions.length > 0) {
      const receiptsBatch = new client.BatchRequest()
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
        .catch(e => {
          Logger.error(`Receipt transaction fetching failed`)
          throw e
        })
      receiptsBatch.execute()
    }
  } catch (e) {
    Logger.error(`error fetching the block ${JSON.stringify(result)}`)
    throw e
  }
}

module.exports = async ({ blockchain }) => {
  if (endpoint(blockchain) === undefined) {
    // We disable this blockchain throw the env variables but
    // if the env variable is present but empty this is an error
    return null
  }
  if (!endpoint(blockchain)) throw new InvalidBlockchainError(blockchain)

  const client = new Web3(new Web3.providers.WebsocketProvider(endpoint(blockchain)))

  // client.eth.defaultBlock = 'latest'
  const subscription = await client.eth.subscribe('newBlockHeaders', (err, result) => {
    if (err) { Logger.error(err) }
  })
  subscription
    .on('changed', () => Logger.info(`Websocket ${blockchain} changed`))
    .on('error', error => Logger.error(`error ${error}`))
    .on('data', handlerNewBlock(client, blockchain))
}
