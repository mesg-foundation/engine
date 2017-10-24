module.exports = ({ transaction, block }) => ({
  blockId: block.number.toString(),
  fees: transaction.gasUsed.toString(),
  from: transaction.from,
  payload: {},
  to: transaction.to,
  transactionId: transaction.hash,
  value: transaction.value,
  executedAt: new Date(block.timestamp * 1000)
})
