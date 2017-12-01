const uniqueKeyGenerator = require('../../uniqueKeyGenerator')

// Generate a unique sortable key for Ethereum events. We concidere that an Ethereum event is
// identified by the following data :
//  - type : the type of blockchain, here always "ETHEREUM"
//  - blockchain : the chain for the blockchain, "TESTNET", "MAINNET"
//  - block : Number of the block that process the event
//  - transaction : Transaction index in the block (cannot take the hash because one event that happen later can have a smaller hash)
//  - log : Log index in the transaction
//
// This will generate somethid with the following exemple data
// ETHEREUM,TESTNET,243134,12,2
// This set of data will ensure that if the event is comming another time we will be able to detect if it has already been processed or no

module.exports = ({ type, blockchain, block, transaction, log }) => uniqueKeyGenerator([
  type,
  blockchain,
  block.number,
  transaction.transactionIndex,
  log ? log.logIndex : null
])
