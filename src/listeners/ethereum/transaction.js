const EthereumConnector = require('../../connectors/ethereum')

const match = trigger => false // TODO trigger.contract && trigger.eventName

const createListener = trigger => {
  const connector = EthereumConnector(trigger.contract.chain)
  const listener = connector
    .filter({
      fromBlock: 'latest',
      toBlock: 'latest',
      address: [trigger.contract.address]
    })
  return {
    watch: callback => listener.watch((error, data) => error
      ? callback(error)
      : Promise.all([
        Promise.resolve(data),
        connector.getTransactionReceiptPromise(data.transactionHash),
        connector.getTransactionPromise(data.transactionHash)
      ])
        .then(results => results.reduce((acc, value) => Object.assign(acc, value), {}))
        .then(data => callback(null, data))
        .catch(error => callback(error))
    ),
    stopWatching: listener.stopWatching
  }
}

module.exports = { match, createListener }