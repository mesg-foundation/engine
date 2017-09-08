const ws = require('ws')
const { ApolloClient, createNetworkInterface } = require('apollo-client')
const { SubscriptionClient, addGraphQLSubscriptions } = require('subscriptions-transport-ws')
const queryFetchAll = require('./queries/fetchAll')
const queryUpdated = require('./queries/updated')
const queryDeleted = require('./queries/deleted')
const mutationEvent = require('./queries/mutationEvent')

const headers = {
  Authorization: [
    'Bearer',
    process.env.AUTH_TOKEN
  ].join(' ')
}

const client = new ApolloClient({
  networkInterface: addGraphQLSubscriptions(
    createNetworkInterface({
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

const variables = {
  chain: process.env.CHAIN
}

const fetchAll = () => client
  .query({ query: queryFetchAll, variables })
  .then(x => x.data.allTriggers)

const onDataUpdated = callback => client
  .subscribe({ query: queryUpdated, variables })
  .subscribe({
    next: value => callback(null, value.Trigger.node),
    error: error => callback(error, null)
  })

const onDataDeleted = callback => client
  .subscribe({ query: queryDeleted, variables })
  .subscribe({
    next: value => callback(null, value.Trigger.previousValues.id),
    error: error => callback(error, null)
  })

const writeEvent = (event, trigger) => client
  .mutate({
    mutation: mutationEvent,
    variables: {
      payload: event.args,
      transactionId: event.transactionHash,
      triggerId: trigger.id
    }
  })
  .then(x => x.data.createEvent)

module.exports = {
  fetchAll,
  onDataUpdated,
  onDataDeleted,
  writeEvent
}
