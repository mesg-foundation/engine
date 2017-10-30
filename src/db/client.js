const ws = require('ws')
const { ApolloClient, createBatchingNetworkInterface } = require('apollo-client')
const { SubscriptionClient, addGraphQLSubscriptions } = require('subscriptions-transport-ws')

const headers = {
  Authorization: [
    'Bearer',
    process.env.AUTH_TOKEN
  ].join(' ')
}

module.exports = new ApolloClient({
  queryDeduplication: true,
  networkInterface: addGraphQLSubscriptions(
    createBatchingNetworkInterface({
      uri: process.env.GRAPHQL_HTTP_ENDPOINT,
      opts: {
        headers
      }
    }),
    new SubscriptionClient(process.env.GRAPHQL_WS_ENDPOINT, {
      reconnect: true,
      timeout: 30000,
      connectionParams: headers
    }, ws)
  )
})
