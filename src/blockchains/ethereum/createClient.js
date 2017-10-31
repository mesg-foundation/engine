const Web3 = require('web3')
const name = require('./name')

module.exports = blockchain => {
  const endpoint = process.env[`${name}_${blockchain.toUpperCase()}`]
  if (endpoint === undefined) {
    // We disable this blockchain throw the env variables but
    // if the env variable is present but empty this is an error
    return null
  }
  if (!endpoint) throw new Error(`The endpoint is empty, please set the env variable ${name}_${blockchain.toUpperCase()}`)

  return new Web3(new Web3.providers.WebsocketProvider(endpoint))
}
