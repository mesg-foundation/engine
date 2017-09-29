const gql = require('graphql-tag')

module.exports = gql`
  query($chain: CHAIN!) {
    allTriggers(filter: {
      enable: true,
      connector: {
        ethereumContract: {
          contract: {
            chain: $chain
          }
        }
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
      }
    }
  }`
