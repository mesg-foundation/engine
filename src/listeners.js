const Store = require('./store')
const Db = require('./db')

const blockchainClients = async () => [
  // await require('./blockchains/ethereum')('MAINNET'),
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
        .forEach(({ trigger, event }) => Db.writeEvent(trigger, event))
    })
  })
}

module.exports = {
  start
}