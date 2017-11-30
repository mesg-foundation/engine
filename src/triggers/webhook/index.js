const uniqueKeyGenerator = require('../../uniqueKeyGenerator')

module.exports = trigger => {
  const { key } = trigger.connector.webhook

  return {
    match: ({ type, request }) => {
      if (type !== 'HTTP') { return false }
      return request.params.key === key
    },
    normalizeEvent: ({ type, request, app }) => ({
      key: uniqueKeyGenerator([
        type,
        request.id
      ]),
      blockId: app.name,
      fees: '0',
      from: request.ip,
      payload: request.body,
      to: '',
      transactionId: request.id,
      value: '',
      executedAt: new Date()
    })
  }
}
