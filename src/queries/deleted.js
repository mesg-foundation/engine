const gql = require('graphql-tag')

module.exports = gql`
  subscription($chain: CHAIN!) {
    Trigger(
      filter: {
        mutation_in: [DELETED]
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
      previousValues {
        id
      }
    }
  }`
