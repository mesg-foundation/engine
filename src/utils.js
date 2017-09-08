const getContract = (web3, trigger) => web3.eth
  .contract(trigger.contract.abi)
  .at(trigger.contract.address)

const isValidEvent = (contract, eventName) => !!contract[eventName]

const createListenerForEvent = (contract, eventName) => isValidEvent(contract, eventName) && contract[eventName](null, {
  fromBlock: 'latest',
  toBlock: 'latest'
})

module.exports = {
  getContract,
  isValidEvent,
  createListenerForEvent
}
