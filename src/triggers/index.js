const { InvalidTrigger } = require('../errors')

const create = {
  ETHEREUM_TRANSACTION: require('./ethereum/transaction'),
  ETHEREUM_CONTRACT: require('./ethereum/contract')
}

module.exports = trigger => {
  const createFunction = create[trigger.connector.connectorType]
  if (!createFunction) throw new InvalidTrigger(`cannot find any factory for ${trigger.connector.connectorType} connector`)
  return {
    id: trigger.id,
    ...createFunction(trigger)
  }
}