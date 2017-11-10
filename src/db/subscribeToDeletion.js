const Logger = require('../logger')
const Store = require('../store')
const { createSubscription } = require('./client')
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
  createSubscription(query, value => Store.remove(value.Trigger.previousValues.id))
}
