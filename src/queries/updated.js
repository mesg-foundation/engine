const gql = require('graphql-tag')

module.exports = gql`
  subscription($chain: CHAIN!) {
    Trigger(
      filter: {
        mutation_in: [CREATED, UPDATED],
        node: {
          contract: {
            chain: $chain
          }
        }
      }
    ) {
      node {
        contract {
          abi
          address
        }
        id
        enable
        eventName
      }
    }
  }`