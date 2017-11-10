const create = {
  ETHEREUM_TRANSACTION: require('./ethereum/transaction'),
  ETHEREUM_CONTRACT: require('./ethereum/contract'),
  ETHEREUM_TOKEN: require('./ethereum/contract')
}

module.exports = trigger => {
  const createFunction = create[trigger.connector.connectorType]
  if (!createFunction) throw new Error(`${trigger.connector.connectorType} is not a valid connector: ${Object.keys(create).join(', ')}`)
  return {
    id: trigger.id,
    ...createFunction(trigger)
  }
}
