const Abi = require('web3-eth-abi')
const Logger = require('../../logger')
const normalizeEvent = require('./normalizeEvent')
const generateKey = require('./generateKey');

module.exports = trigger => {
  const { eventName, contract } = trigger.connector.ethereumContract || trigger.connector.ethereumToken
  const { chain, address } = contract
  const eventAbi = contract.abi
    .filter(x => x.type === 'event')
    .filter(x => x.name === eventName)[0]

  const encodedEvent = Abi.encodeEventSignature(eventAbi)
  const matchLog = log => (log.topics || []).some(topic => topic === encodedEvent)

  return {
    match: ({ type, blockchain, transaction, block }) => {
      if (type !== 'ETHEREUM') { return false }
      if (blockchain !== chain) { return false }
      if (address.toLowerCase() !== (transaction.to || '').toLowerCase()) { return false }

      if (!transaction.logs) {
        Logger.error('transaction log not valid', { type, blockchain, transaction, block })
      }
      return (transaction.logs || []).some(matchLog)
    },
    normalizeEvent: event => {
      const normalizedEvent = normalizeEvent(event)
      return event.transaction.logs
        .filter(matchLog)
        .map(log => {
          const data = Abi.decodeLog(eventAbi.inputs, log.data, log.topics)
          return {
            ...normalizedEvent,
            key: generateKey({
              ...event,
              log
            }),
            payload: eventAbi.inputs
              .map(x => x.name)
              .reduce((acc, name) => ({
                ...acc,
                [name]: data[name]
              }), {})
          }
        })
    }
  }
}
