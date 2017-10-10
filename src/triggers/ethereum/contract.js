module.exports = trigger => {
  const { eventName, contract } = trigger.connector.ethereumContract
  const { chain, address } = contract

  return {
    match: ({ type, network, transaction, block }) => {
      if (type !== 'ETHEREUM') { return false }
      if (network !== chain) { return false }
      if (address.toLowerCase() !== (transaction.to || '').toLowerCase()) { return false }
      return false
      // TODO return event from the transaction
    }
  }
}