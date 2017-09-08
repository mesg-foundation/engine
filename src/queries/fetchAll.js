const gql = require('graphql-tag')

module.exports = gql`
  query($chain: CHAIN!) {
    allTriggers(filter: {
      enable: true,
      contract: {
        chain: $chain
      }
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
