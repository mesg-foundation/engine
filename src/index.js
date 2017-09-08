require('dotenv').config()
require('isomorphic-fetch')
const Web3 = require('web3')
const web3 = new Web3(new Web3.providers.HttpProvider(process.env.NODE_ADDRESS))

const Store = require('./store')(web3)
const DB = require('./db')

const handleEvent = ({ event, trigger }) => DB.writeEvent(event, trigger)

const executeIfNoError = callback => (err, data) => err
  ? console.error(err) // TODO need to handle error
  : callback(data)

DB.fetchAll()
  .then(triggers => triggers
    .map(trigger => Store.add(trigger, executeIfNoError(handleEvent)))
  )

DB.onDataUpdated(executeIfNoError(trigger => trigger.enable
  ? Store.add(trigger, executeIfNoError(handleEvent))
  : Store.remove(trigger.id)))

DB.onDataDeleted(executeIfNoError(triggerId => Store.remove(triggerId)))
