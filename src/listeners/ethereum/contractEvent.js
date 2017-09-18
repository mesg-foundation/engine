const EthereumConnector = require('../../connectors/ethereum')

module.exports = trigger => {
  const onEvent = EthereumConnector(trigger.contract.chain)
    .contract(trigger.contract.abi)
    .at(trigger.contract.address)[trigger.eventName]
  return onEvent
    ? onEvent(null, {
      fromBlock: 'latest',
      toBlock: 'latest'
    })
    : null
}