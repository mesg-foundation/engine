const gql = require('graphql-tag')

module.exports = gql`
  subscription($chain: CHAIN!) {
    Trigger(
      filter: {
        mutation_in: [DELETED]
        node: {
          contract: {
            chain: $chain
          }
        }
      }
    ) {
      previousValues {
        id
      }
    }
  }`
