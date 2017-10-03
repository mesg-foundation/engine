const createClient = require('./client')

const match = trigger => trigger.connector.connectorType === 'ETHEREUM_TRANSACTION'

const createListener = async trigger => {
  const { address, chain } = trigger.connector.ethereumTransaction
  const client = await createClient(chain)
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
