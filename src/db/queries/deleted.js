const gql = require('graphql-tag')

module.exports = gql`
  subscription {
    Trigger(
      filter: {
        mutation_in: [DELETED]
      }
    ) {
      previousValues {
        id
      }
    }
  }`
