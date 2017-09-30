const gql = require('graphql-tag')

module.exports = gql`
  subscription($chain: ETHEREUM_BLOCKCHAIN!) {
    Trigger(
      filter: {
        mutation_in: [DELETED]
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
      previousValues {
        id
      }
    }
  }`
