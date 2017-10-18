const normalizeEvent = require('./normalizeEvent')

module.exports = trigger => {
  const { chain, address } = trigger.connector.ethereumTransaction

  return {
    match: ({ type, blockchain, transaction, block }) => {
      if (type !== 'ETHEREUM') { return false }
      if (blockchain !== chain) { return false }
      return address.toLowerCase() === (transaction.from || '').toLowerCase() ||
             address.toLowerCase() === (transaction.to || '').toLowerCase()
    },
    normalizeEvent
  }
}
