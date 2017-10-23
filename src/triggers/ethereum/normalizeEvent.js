module.exports = ({ type, blockchain, block, transaction }) => ({
  keys: [
    type,
    blockchain,
    transaction.hash
  ],
  blockId: block.number.toString(),
  fees: transaction.gasUsed.toString(),
  from: transaction.from,
  payload: {},
  to: transaction.to,
  transactionId: transaction.hash,
  value: transaction.value,
  executedAt: new Date(block.timestamp * 1000)
})
