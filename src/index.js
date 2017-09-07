require('dotenv').config()
require('isomorphic-fetch')
const Web3 = require('web3')
const web3 = new Web3(new Web3.providers.HttpProvider(process.env.NODE_ADDRESS))

const Store = require('./store')(web3)
const DB = require('./db')

const handleEvent = (err, { event, trigger }) => err
  ? console.log(err) // TODO need to handle this error
  : DB.writeEvent(event, trigger)

DB.fetchAll()
  .then(triggers => triggers
    .map(trigger => Store.add(trigger, handleEvent))
  )

DB.onDataCreated((err, trigger) => Store.add(trigger, handleEvent))
DB.onDataUpdated((err, trigger) => Store.add(trigger, handleEvent))
DB.onDataDeleted((err, triggerId) => Store.remove(triggerId))