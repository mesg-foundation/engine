const Logger = require('../logger')
const Store = require('../store')
const client = require('./client')
const gql = require('graphql-tag')

const query = gql`subscription {
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

module.exports = () => {
  Logger.info('Connecting to trigger delete...')
  client
    .subscribe({ query })
    .subscribe({
      next: value => Store.remove(value.Trigger.previousValues.id),
      error: error => Logger.error('Graphcool delete subscription error', { error })
    })
}
