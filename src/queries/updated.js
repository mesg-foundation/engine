const gql = require('graphql-tag')

module.exports = gql`
  subscription($chain: CHAIN!) {
    Trigger(
      filter: {
        mutation_in: [CREATED, UPDATED],
        node: {
          connector: {
            ethereumContract: {
              contract: {
                chain: $chain
              }
            }
          }
        }
      }
    ) {
      node {
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
    }
  }`
