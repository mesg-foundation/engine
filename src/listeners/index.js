const blockchainClients = async () => [
  await require('./blockchainClients/ethereum')('KOVAN')
]

const start = async () => {
  console.debug('initializing all blockchains connections')
  const clients = await blockchainClients()
  console.debug('listening for transactions')
  clients.forEach(({ type, network, onTransaction }) => {
    onTransaction((transaction, block) => {
      console.log(transaction.from)
    })
  })
}

module.exports = {
  start
}