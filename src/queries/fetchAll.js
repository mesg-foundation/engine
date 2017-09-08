const gql = require('graphql-tag')

module.exports = gql`
  query {
    allTriggers(filter: {
      enable: true
    }) {
      contract {
        abi
        address
      }
      id
      enable
      eventName
    }
  }`