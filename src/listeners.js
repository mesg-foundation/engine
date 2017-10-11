const Store = require('./store')
const Db = require('./db')

const blockchainClients = async () => [
  // await require('./blockchains/ethereum')('MAINNET')
  await require('./blockchains/ethereum')('KOVAN')
]

const start = async () => {
  console.debug('initializing all blockchains connections')
  const clients = await blockchainClients()
  console.debug('listening for transactions')
  clients.forEach(({ type, network, onTransaction }) => {
    onTransaction((transaction, block) => {
      Store
        .matchingTriggers({ type, network, transaction, block })
        .map(trigger => ({
          trigger,
          event: trigger.normalizeEvent({ transaction, block })
        }))
        .forEach(result => Array.isArray(result.event)
          ? result.event.map(event => Db.writeEvent(result.trigger, event))
          : Db.writeEvent(result.trigger, result.event))
    })
  })
}

module.exports = {
  start
}