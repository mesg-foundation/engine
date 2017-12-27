const Logger = require('../../logger')

module.exports = async (client, number) => {
  try {
    let block = null
    let retry = 5
    while (!block && retry > 0) {
      block = await client.eth.getBlock(number, true)
      retry = retry - 1
    }
    if (!block) { throw new Error('Max retry') }
    return block
  } catch (e) {
    Logger.error('error fetching the block', { number })
    throw e
  }
}
