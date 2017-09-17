require('dotenv').config()
require('isomorphic-fetch')
require('newrelic')
const bugsnag = require('bugsnag')
const Web3 = require('web3')

bugsnag.register(process.env.BUGSNAG_KEY)

const web3 = new Web3(new Web3.providers.HttpProvider(process.env.NODE_ADDRESS))

const Store = require('./store')(web3)
const DB = require('./db')

const handleEvent = ({ event, trigger }) => DB.writeEvent(event, trigger)
  .catch(bugsnag.notify)

const executeIfNoError = callback => (err, data) => {
  if (err) {
    bugsnag.notify(new Error(err))
    console.error(err)
  } else {
    callback(data)
  }
}

DB.fetchAll()
  .then(triggers => triggers
    .map(trigger => Store.add(trigger, executeIfNoError(handleEvent)))
  )

DB.onDataUpdated(executeIfNoError(trigger => trigger.enable
  ? Store.add(trigger, executeIfNoError(handleEvent))
  : Store.remove(trigger.id)))

DB.onDataDeleted(executeIfNoError(triggerId => Store.remove(triggerId)))
