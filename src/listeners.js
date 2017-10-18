const Logger = require('./logger')
const Store = require('./store')
const Db = require('./db')

const blockchainClients = async () => [
  // await require('./blockchains/ethereum')('MAINNET'),
  await require('./blockchains/ethereum')('KOVAN')
]

const handleTransaction = (type, network) => (transaction, block) => Store
  .matchingTriggers({ type, network, transaction, block })
  .map(trigger => ({
    trigger,
    event: trigger.normalizeEvent({ transaction, block })
  }))
  .forEach(result => Array.isArray(result.event)
    ? result.event.map(event => Db.writeEvent(result.trigger, event))
    : Db.writeEvent(result.trigger, result.event))

const connectClientToTransactions = ({ type, network, onTransaction }) =>
  onTransaction(handleTransaction(type, network))

const start = async () => {
  Logger.info('initializing all blockchains connections')
  const clients = await blockchainClients()
  Logger.info('listening for transactions')
  clients.forEach(connectClientToTransactions)
}

module.exports = {
  start
}
