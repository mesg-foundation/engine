const Logger = require('../logger')
const Store = require('../store')
const client = require('./client')
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
        }
      }
    }
  }
}`

module.exports = () => {
  Logger.info('Connecting to trigger update...')
  client
    .subscribe({ query })
    .subscribe({
      next: value => Store.update(value.Trigger.node),
      error: error => Logger.error('Graphcool update/create subscription error', { error })
    })
}
