const gql = require('graphql-tag')

module.exports = gql`
  query($chain: ETHEREUM_BLOCKCHAIN!) {
    allTriggers(filter: {
      enable: true,
      connector: {
        OR: [
          {
            ethereumContract: {
              contract: {
                chain: $chain
              }
            }
          },
          {
            ethereumTransaction: {
              chain: $chain
            }
          }
        ]
      }
    }) {
      id
      enable
      connector {
        connectorType
        ethereumContract {
          eventName
          contract {
            abi
            address
          }
        }
        ethereumTransaction {
          address
          chain
        }
      }
    }
  }`
