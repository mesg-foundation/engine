const gql = require('graphql-tag')

module.exports = gql`
  query {
    allTriggers(filter: {
      enable: true
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
