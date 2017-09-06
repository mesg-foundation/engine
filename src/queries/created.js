const gql = require('graphql-tag')

module.exports = gql`
  subscription {
    Trigger(
      filter: {
        mutation_in: [CREATED]
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