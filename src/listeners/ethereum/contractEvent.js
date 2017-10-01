const createClient = require('./client')

const match = trigger => trigger.connector.connectorType === 'ETHEREUM_CONTRACT'

const createListener = trigger => {
  const { contract, eventName } = trigger.connector.ethereumContract
  const client = createClient(contract.chain)
  const onEvent = client
    .contract(contract.abi)
    .at(contract.address)[eventName]
  if (!onEvent) { return null }
  const listener = onEvent(null, {
    fromBlock: 'latest',
    toBlock: 'latest'
  })
  return {
    watch: callback => listener.watch(client.handleEvent(callback)),
    stopWatching: listener.stopWatching
  }
}

module.exports = { match, createListener }
