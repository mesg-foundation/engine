const Store = require('./store')

const blockchainClients = async () => [
  await require('./blockchains/ethereum')('KOVAN'),
  await require('./blockchains/ethereum')('MAINNET')
]

const start = async () => {
  console.debug('initializing all blockchains connections')
  const clients = await blockchainClients()
  console.debug('listening for transactions')
  clients.forEach(({ type, network, onTransaction }) => {
    onTransaction((transaction, block) => {
      console.log(`transaction ${network}`)
      Store
        .matchingTriggers({ type, network, transaction, block })
        .forEach(trigger => trigger.emitEvent(transaction, block))
    })
  })
}

module.exports = {
  start
}