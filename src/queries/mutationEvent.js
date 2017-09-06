const gql = require('graphql-tag')

module.exports = gql`
  mutation(
    $payload: Json!,
    $transactionId: String!,
    $triggerId: ID!
  ) {
    createEvent(
      payload: $payload
      transactionId: $transactionId
      triggerId: $triggerId
    ) {
      id
    }
  }`