const { emitRawBlock, emitRawTransaction } = require('../../eventEmitter')
const Block = require('../../db/block')
const type = require('./name')
const fetchBlock = require('./fetchBlock')
const transactionsWithReceiptBatch = require('./transactionsWithReceiptBatch')

module.exports = (client, blockchain) => async blockNumber => {
  const block = await fetchBlock(client, blockNumber)
  emitRawBlock({ type, blockchain, block })
  const transactions = await transactionsWithReceiptBatch(client, block.transactions)
  transactions.forEach(transaction => {
    emitRawTransaction({ type, blockchain, block, transaction })
  })
  await Block.processed({ type, blockchain }, blockNumber)
}
