const { testConnection } = require('../../utils')

module.exports = async ({
  type,
  network,
  onTransaction,
  isConnected
}) => {
  if (!type) throw new Error('type is missing')
  if (!network) throw new Error('network is missing')
  if (!onTransaction) throw new Error('onTransaction is missing')
  if (!isConnected) throw new Error('isConnected is missing')

  await testConnection(isConnected, `${type}/${network}`)

  return {
    type,
    network,
    onTransaction
  }
}