const { InvalidClientError } = require('../errors')
const { testConnection } = require('../utils')

module.exports = async ({
  type,
  network,
  onTransaction,
  isConnected
}) => {
  if (!type) throw new InvalidClientError('type is missing')
  if (!network) throw new InvalidClientError('network is missing')
  if (!onTransaction) throw new InvalidClientError('onTransaction is missing')
  if (!isConnected) throw new InvalidClientError('isConnected is missing')

  await testConnection(isConnected, `${type}/${network}`)

  return {
    type,
    network,
    onTransaction
  }
}