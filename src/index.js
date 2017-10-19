require('dotenv').config()
require('isomorphic-fetch')
require('newrelic')
require('bugsnag').register(process.env.BUGSNAG_KEY)

const eventEmitter = require('./eventEmitter')
const Logger = require('./logger')
const DB = require('./db')
const Store = require('./store')
const initializeBlockchains = require('./blockchains')

const handleRawTransaction = transactionArgs => Store
  .all()
  .filter(trigger => trigger.match(transactionArgs))
  .map(trigger => ({
    trigger,
    event: trigger.normalizeEvent(transactionArgs)
  }))
  .forEach(result => Array.isArray(result.event)
    ? result.event.map(event => DB.writeEvent(result.trigger, event))
    : DB.writeEvent(result.trigger, result.event))

const startApp = async () => {
  eventEmitter.create()
  eventEmitter.on('RAW_TRANSACTION', handleRawTransaction)

  Logger.info('init database')
  await DB.init()
  
  Logger.info('initializing all blockchains connections')
  await initializeBlockchains()
}

startApp()
