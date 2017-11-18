const normalizeEvent = require('./normalizeEvent')

module.exports = trigger => {
  const { chain, address, matchType } = trigger.connector.ethereumTransaction
  const addr = address.toLowerCase()
  const matchFn = {
    ANY: tx => addr === (tx.from || '').toLowerCase() || addr === (tx.to || '').toLowerCase(),
    FROM: tx => addr === (tx.from || '').toLowerCase(),
    TO: tx => addr === (tx.to || '').toLowerCase()
  }[matchType || 'ANY']

  return {
    match: ({ type, blockchain, transaction, block }) => {
      if (type !== 'ETHEREUM') { return false }
      if (blockchain !== chain) { return false }
      return matchFn(transaction)
    },
    normalizeEvent
  }
}
