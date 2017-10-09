const gql = require('graphql-tag')

module.exports = gql`
  mutation(
    $blockId: String!,
    $fees: String!,
    $from: String!,
    $payload: Json!,
    $to: String!
    $transactionId: String!
    $value: String!,
    $triggerId: ID!
  ) {
    createEvent(
      blockId: $blockId,
      fees: $fees,
      from: $from,
      payload: $payload,
      to: $to,
      transactionId: $transactionId,
      value: $value,
      triggerId: $triggerId,
    ) {
      id
    }
  }`
