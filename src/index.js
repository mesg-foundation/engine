require('dotenv').config()
require('isomorphic-fetch')
require('newrelic')
const Bugsnag = require('bugsnag')
Bugsnag.register(process.env.BUGSNAG_KEY)

const eventEmitter = require('./eventEmitter')
const Logger = require('./logger')
const DB = require('./db')
const Store = require('./store')
const initializeBlockchains = require('./blockchains')

const handleRawTransaction = transactionArgs => {
  const events = Store.all()
    .filter(trigger => trigger.match(transactionArgs))
    .map(trigger => ({
      trigger,
      event: trigger.normalizeEvent(transactionArgs)
    }))
    .reduce((list, result) => {
      return [
        ...list,
        ...Array.isArray(result.event)
          ? result.event.map(event => ({ trigger: result.trigger, event }))
          : [result]
      ]
    }, [])
  if (events.length > 0) {
    Logger.info('Submit events', { events })
  }
  events
    .sort((a, b) => a.event.key.localeCompare(b.event.key))
    .map(DB.createEvent)
}

const startApp = async () => {
  try {
    eventEmitter.create()
    eventEmitter.on('RAW_TRANSACTION', handleRawTransaction)
    eventEmitter.on('RAW_BLOCK', ({ type, blockchain, block }) => Logger.info('New block', { type, blockchain, block: block.number }))

    Logger.info('init database')
    await DB.init()

    Logger.info('initializing all blockchains connections')
    await initializeBlockchains()
  } catch (e) {
    Bugsnag.notify(e)
    console.log(e)
    Logger.error(e)
    process.exit(-1)
  }
}

startApp()
