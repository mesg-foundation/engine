const Web3 = require('web3')

const connectors = {}

module.exports = chain => {
  if (!connectors[chain]) {
    connectors[chain] = (new Web3(new Web3.providers.HttpProvider(process.env.NODE_ADDRESS))).eth
  }
  return connectors[chain]
}