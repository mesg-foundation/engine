const Logger = require('../logger')
const { client } = require('./client')
const { KEY_SEPARATOR } = require('../uniqueKeyGenerator')
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
  $key: String!
) {
  createEvent(
    key: $key,
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
      key: [trigger.id, event.key].join(KEY_SEPARATOR),
      triggerId: trigger.id
    }
  })
  .then(x => x.data.createEvent)
  .catch(e => {
    if (e.graphQLErrors.every(x => x.code === 3010)) { // duplicate event
      return;
    }
    console.log(e)
    Logger.error('Event creation fails', { event, trigger })
    throw e
  })
