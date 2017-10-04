const { InvalidEventError } = require('../../errors')
const createClient = require('./client')

const match = trigger => trigger.connector.connectorType === 'ETHEREUM_CONTRACT'

const createListener = async trigger => {
  const { contract, eventName } = trigger.connector.ethereumContract
  const client = await createClient(contract.chain)
  const onEvent = client.eth
    .contract(contract.abi)
    .at(contract.address)[eventName]
  if (!onEvent) { throw new InvalidEventError(eventName) }
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
