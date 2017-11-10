const Logger = require('../../logger')

module.exports = async (client, number) => {
  try {
    const block = client.eth.getBlock(number, true)
    return block
  } catch (e) {
    Logger.error('error fetching the block', { number })
    throw e
  }
}
