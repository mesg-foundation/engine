const Logger = require('../logger')
const Store = require('../store')
const { createSubscription } = require('./client')
const gql = require('graphql-tag')

const query = gql`subscription {
  Trigger(
    filter: {
      mutation_in: [CREATED, UPDATED]
    }
  ) {
    node {
      id
      enable
      connector {
        connectorType
        ethereumContract {
          eventName
          contract {
            abi
            address
            chain
          }
        },
        ethereumToken {
          eventName
          contract {
            abi
            address
            chain
          }
        },
        ethereumTransaction {
          chain
          address
          matchType
        }
        webhook {
          key
        }
      }
    }
  }
}`

module.exports = () => {
  Logger.info('Connecting to trigger update...')
  createSubscription(query, value => Store.update(value.Trigger.node))
}
