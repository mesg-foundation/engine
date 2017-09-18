const ethereum = {
  contractEvent: require('./ethereum/contractEvent')
}

const isEthereumContractEvent = x => x.contract && x.eventName

module.exports = trigger => {
  if (isEthereumContractEvent(trigger)) { return ethereum.contractEvent(trigger) }
  throw new Error(`${trigger.id} does not have any valid listener`)
}