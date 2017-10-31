const Web3 = require('web3')
const { InvalidBlockchainError } = require('../../errors')
const name = require('./name')

module.exports = blockchain => {
  const endpoint = process.env[`${name}_${blockchain.toUpperCase()}`]
  if (endpoint === undefined) {
    // We disable this blockchain throw the env variables but
    // if the env variable is present but empty this is an error
    return null
  }
  if (!endpoint) throw new InvalidBlockchainError(blockchain)

  return new Web3(new Web3.providers.WebsocketProvider(endpoint))
}