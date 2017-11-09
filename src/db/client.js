const ws = require('ws')
const { ApolloClient, createBatchingNetworkInterface } = require('apollo-client')
const { SubscriptionClient, addGraphQLSubscriptions } = require('subscriptions-transport-ws')
const Logger = require('../logger')

const headers = {
  Authorization: [
    'Bearer',
    process.env.AUTH_TOKEN
  ].join(' ')
}

const client = new ApolloClient({
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

const createSubscription = (query, next) => {
  Logger.info('(re)subscribe')
  const subscription = client
    .subscribe({ query })
    .subscribe({
      next,
      error: error => {
        Logger.error('Subscription error', { error })
        subscription.unsubscribe()
        createSubscription(query, next)
      }
    })
  return subscription
}

module.exports = {
  client,
  createSubscription
}