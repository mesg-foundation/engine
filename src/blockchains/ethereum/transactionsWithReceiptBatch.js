const Logger = require('../../Logger')

const receiptBatcher = (client, batch) => transaction => new Promise((resolve, reject) => {
  batch.add(client.eth.getTransactionReceipt.request(
    transaction.hash,
    (err, receipt) => err
      ? reject(err)
      : resolve({ ...transaction, ...receipt })
  ))
})

module.exports = async (client, transactions) => {
  if (transactions.length <= 0) { return Promise.resolve([]) }
  try {
    const batch = new client.BatchRequest()
    const addFetchToBatch = receiptBatcher(client, batch)
    const promises = transactions.map(addFetchToBatch)
    batch.execute()
    const results = await Promise.all(promises)
    return results
  } catch (e) {
    Logger.error('Receipt transaction fetching failed', { transactions: transactions.map(x => x.hash) })
    throw e
  }
}
