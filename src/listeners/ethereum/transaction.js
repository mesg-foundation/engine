const EthereumConnector = require('../../connectors/ethereum')

const match = trigger => trigger.connector.connectorType === 'ETHEREUM_TRANSACTION'

const createListener = trigger => {
  const { address, chain } = trigger.connector.ethereumTransaction
  const ethConnector = EthereumConnector(chain)
  const listener = ethConnector
    .filter({
      fromBlock: 'latest',
      toBlock: 'latest',
      address: [address]
    })
  return {
    watch: callback => listener.watch(ethConnector.handleEvent(callback)),
    stopWatching: listener.stopWatching
  }
}

module.exports = { match, createListener }
