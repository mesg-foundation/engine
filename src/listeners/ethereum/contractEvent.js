const EthereumConnector = require('../../connectors/ethereum')

const match = trigger => trigger.connector.connectorType === 'ETHEREUM_CONTRACT'

const createListener = trigger => {
  const { contract, eventName } = trigger.connector.ethereumContract
  const ethConnector = EthereumConnector(contract.chain)
  const onEvent = ethConnector
    .contract(contract.abi)
    .at(contract.address)[eventName]
  if (!onEvent) { return null }
  const listener = onEvent(null, {
    fromBlock: 'latest',
    toBlock: 'latest'
  })
  return {
    watch: callback => listener.watch(ethConnector.handleEvent(callback)),
    stopWatching: listener.stopWatching
  }
}

module.exports = { match, createListener }
