const EthereumConnector = require('../../connectors/ethereum')

const match = trigger => trigger.connector.connectorType === 'ETHEREUM_CONTRACT'

const parseResult = data => ({
  blockId: data.blockNumber.toString(),
  fees: data.gasUsed.toString(),
  from: data.from,
  payload: data.args,
  to: data.to,
  transactionId: data.transactionHash,
  value: data.value
})

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
    watch: callback => listener.watch((error, data) => error
      ? callback(error)
      : Promise.all([
        Promise.resolve(data),
        ethConnector.getTransactionReceiptPromise(data.transactionHash),
        ethConnector.getTransactionPromise(data.transactionHash)
      ])
        .then(results => results.reduce((acc, value) => Object.assign(acc, value), {}))
        .then(parseResult)
        .then(data => callback(null, data))
        .catch(error => callback(error))
    ),
    stopWatching: listener.stopWatching
  }
}

module.exports = { match, createListener }