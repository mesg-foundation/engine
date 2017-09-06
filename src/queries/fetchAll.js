const gql = require('graphql-tag')

module.exports = gql`
  query {
    allTriggers {
      contract {
        abi
        address
      }
      id
      eventName
    }
  }`