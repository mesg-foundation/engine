const SolidityEvent = require('web3/lib/web3/event')
const Logger = require('../../logger')
const normalizeEvent = require('./normalizeEvent')

const matchLogFromTopics = topics => log => (log.topics || [])
  .some(topic => topics.indexOf(topic) >= 0)

module.exports = trigger => {
  const { eventName, contract } = trigger.connector.ethereumContract || trigger.connector.ethereumToken
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

      if (!transaction.logs) {
        Logger.error(`transaction log not valid ${transaction.hash}`)
      }
      return (transaction.logs || []).some(matchLog)
    },
    normalizeEvent: event => {
      const normalizedEvent = normalizeEvent(event)
      return event.transaction.logs
        .filter(matchLog)
        .map(log => ({
          ...normalizedEvent,
          // we need to copy the log because web3 modify it directly and remove the topic field
          // needed to know if the contract match the event so if 2+ triggers listen the same contract
          // event then only the first one would be notified if we don't pass a copy to the decode
          // function of web3
          payload: solidityEvent.decode({ ...log }).args
        }))
    }
  }
}
