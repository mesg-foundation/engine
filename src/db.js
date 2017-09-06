const ws = require('ws')
const { ApolloClient, createNetworkInterface } = require('apollo-client')
const { SubscriptionClient, addGraphQLSubscriptions } = require('subscriptions-transport-ws')
const queryFetchAll = require('./queries/fetchAll')
const queryCreated = require('./queries/created')
const queryUpdated = require('./queries/updated')
const queryDeleted = require('./queries/deleted')
const mutationEvent = require('./queries/mutationEvent')

const headers = {
  Authorization: [
    'Bearer',
    'eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpYXQiOjE0OTk1OTA0OTEsImNsaWVudElkIjoiY2o0bnB3dW9mb3FzZDAxMThjcG1zaDdwMSIsInByb2plY3RJZCI6ImNqNG5wd3VvZm9xc2MwMTE4b2VsczdsaW8iLCJwZXJtYW5lbnRBdXRoVG9rZW5JZCI6ImNqNHdoZ2s4dmV0MzEwMTMxeGdla3FsNzcifQ.VXyiJ0ZNLiZDdF0cjGvx0zIYksCtDBZz-u2kLHHIbkQ'
  ].join(' ')
}

const client = new ApolloClient({
  networkInterface: addGraphQLSubscriptions(
    createNetworkInterface({
      uri: 'https://api.graph.cool/simple/v1/eth-hook',
      opts: {
        headers,
      }
    }),
    new SubscriptionClient('wss://subscriptions.graph.cool/v1/eth-hook', {
      reconnect: true,
      connectionParams: headers,
    }, ws)
  )
})

const fetchAll = () => client
  .query({ query: queryFetchAll })
  .then(x => x.data.allTriggers)

const onDataCreated = callback => client
  .subscribe({ query: queryCreated })
  .subscribe({
    next: value => callback(null, value.Trigger.node),
    error: error => callback(error, null)
  })

const onDataUpdated = callback => client
  .subscribe({ query: queryUpdated })
  .subscribe({
    next: value => callback(null, value.Trigger.node),
    error: error => callback(error, null)
  })

const onDataDeleted = callback => client
  .subscribe({ query: queryDeleted })
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
  onDataCreated,
  onDataUpdated,
  onDataDeleted,
  writeEvent
}