const gql = require('graphql-tag')

module.exports = gql`
  subscription($chain: ETHEREUM_BLOCKCHAIN!) {
    Trigger(
      filter: {
        mutation_in: [CREATED, UPDATED],
        node: {
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
          },
          ethereumTransaction {
            chain
            address
          }
        }
      }
    }
  }`
