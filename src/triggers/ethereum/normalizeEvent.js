const generateKey = require('./generateKey');

module.exports = event => {
  const { type, blockchain, block, transaction } = event
  return {
    key: generateKey(event),
    blockId: block.number.toString(),
    fees: transaction.gasUsed.toString(),
    from: transaction.from,
    payload: {},
    to: transaction.to,
    transactionId: transaction.hash,
    value: transaction.value,
    executedAt: new Date(block.timestamp * 1000)
  }
}