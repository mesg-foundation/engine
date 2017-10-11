const gql = require('graphql-tag')

module.exports = gql`
  subscription {
    Trigger(
      filter: {
        mutation_in: [CREATED, UPDATED]
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
              chain
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
