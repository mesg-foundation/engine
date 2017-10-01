const createClient = require('./client')

const match = trigger => trigger.connector.connectorType === 'ETHEREUM_TRANSACTION'

const createListener = trigger => {
  const { address, chain } = trigger.connector.ethereumTransaction
  const client = createClient(chain)
  const listener = client
    .filter({
      fromBlock: 'latest',
      toBlock: 'latest',
      address: [address]
    })
  return {
    watch: callback => listener.watch(client.handleEvent(callback)),
    stopWatching: listener.stopWatching
  }
}

module.exports = { match, createListener }
