const Logger = require('../logger')
const fetchTrigger = require('./fetch')
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
    }
  }
}`

module.exports = () => {
  Logger.info('Connecting to trigger update...')
  createSubscription(query, value => fetchTrigger(value.Trigger.node.id))
}
