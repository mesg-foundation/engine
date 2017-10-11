const DB = require('../../db')

module.exports = trigger => {
  const { chain, address } = trigger.connector.ethereumTransaction

  return {
    match: ({ type, network, transaction, block }) => {
      if (type !== 'ETHEREUM') { return false }
      if (network !== chain) { return false }
      return address.toLowerCase() === (transaction.from || '').toLowerCase()
          || address.toLowerCase() === (transaction.to || '').toLowerCase()
    },
    normalizeEvent: ({ transaction, block }) => ({
      blockId: block.number.toString(),
      fees: transaction.gasUsed.toString(),
      from: transaction.from,
      payload: {},
      to: transaction.to,
      transactionId: transaction.hash,
      value: transaction.value
    })
  }
}