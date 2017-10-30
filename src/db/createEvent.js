const Logger = require('../logger')
const client = require('./client')
const gql = require('graphql-tag')

const mutation = gql`mutation(
  $blockId: String!,
  $fees: String!,
  $from: String!,
  $payload: Json!,
  $to: String!
  $transactionId: String!
  $value: String!,
  $executedAt: DateTime!,
  $triggerId: ID!,
) {
  createEvent(
    blockId: $blockId,
    fees: $fees,
    from: $from,
    payload: $payload,
    to: $to,
    transactionId: $transactionId,
    value: $value,
    executedAt: $executedAt,
    triggerId: $triggerId,
  ) {
    id
  }
}`

module.exports = ({ trigger, event }) => client
  .mutate({
    mutation,
    variables: {
      ...event,
      triggerId: trigger.id
    }
  })
  .then(x => x.data.createEvent)
  .catch(e => {
    Logger.error(e)
    throw new Error('Event creation fails')
  })
  