const create = {
  ETHEREUM_TRANSACTION: require('./ethereum/transaction'),
  ETHEREUM_CONTRACT: require('./ethereum/contract')
}

module.exports = trigger => {
  const createFunction = create[trigger.connector.connectorType]
  if (!createFunction) throw new Error('not valid trigger')
  return {
    id: trigger.id,
    ...createFunction(trigger)
  }
}