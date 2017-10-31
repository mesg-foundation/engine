const Logger = require('../../logger')
const createClient = require('./createClient')
const blockProcessor = require('./blockProcessor')

module.exports = async ({ blockchain }) => {
  const client = createClient(blockchain)
  const processBlock = blockProcessor(client, blockchain)

  // client.eth.defaultBlock = 'latest'
  const subscription = await client.eth.subscribe('newBlockHeaders', (err, result) => {
    if (err) { Logger.error('Error on subscribe', { blockchain, err }) }
  })
  subscription
    .on('changed', () => Logger.info(`Websocket ${blockchain} changed`))
    .on('error', error => Logger.error('error on ethereum subscription', { blockchain, error }))
    .on('data', block => processBlock(block.number))
}
