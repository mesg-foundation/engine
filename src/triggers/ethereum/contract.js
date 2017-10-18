const SolidityEvent = require('web3/lib/web3/event')
const normalizeEvent = require('./normalizeEvent')

const matchLogFromTopics = topics => log => (log.topics || [])
  .some(topic => topics.indexOf(topic) >= 0)

module.exports = trigger => {
  const { eventName, contract } = trigger.connector.ethereumContract
  const { chain, address } = contract

  const eventAbi = contract.abi
    .filter(x => x.type === 'event')
    .filter(x => x.name === eventName)[0]

  const solidityEvent = new SolidityEvent(null, eventAbi, address)
  const matchLog = matchLogFromTopics(solidityEvent.encode().topics)

  return {
    match: ({ type, blockchain, transaction, block }) => {
      if (type !== 'ETHEREUM') { return false }
      if (blockchain !== chain) { return false }
      if (address.toLowerCase() !== (transaction.to || '').toLowerCase()) { return false }

      return transaction.logs.some(matchLog)
    },
    normalizeEvent: event => {
      const normalizedEvent = normalizeEvent(event)
      return event.transaction.logs
        .filter(matchLog)
        .map(log => Object.assign({}, normalizedEvent, {
          payload: solidityEvent.decode(log).args
        }))
    }
  }
}
