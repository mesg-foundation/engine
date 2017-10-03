const { JsonRPCConnectionError } = require('./errors')

const sleep = ms => new Promise(resolve => setTimeout(resolve, ms))

const testConnection = async (check, endpoint) => {
  let retry = process.env.RETRY_COUNT
  const time = new Date()
  while (!check()) {
    if (retry <= 0) {
      throw new JsonRPCConnectionError(endpoint)
    }
    await sleep(1000)
    console.log(`Connection to ${endpoint} invalid... retrying in 1s`)
    retry = retry - 1
  }
  return Promise.resolve({
    retryCount: retry - process.env.RETRY_COUNT,
    ms: +new Date() - time
  })
}

module.exports = {
  sleep,
  testConnection
}
