const gql = require('graphql-tag')

module.exports = gql`
  subscription {
    Trigger(
      filter: {
        mutation_in: [UPDATED]
      }
    ) {
      node {
        contract {
          abi
          address
        }
        id
        eventName
      }
    }
  }`