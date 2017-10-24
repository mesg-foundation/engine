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
