require('dotenv').config()
require('isomorphic-fetch')
require('newrelic')
const bugsnag = require('bugsnag')
bugsnag.register(process.env.BUGSNAG_KEY)

const Store = require('./store')
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
