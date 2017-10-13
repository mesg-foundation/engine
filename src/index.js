require('dotenv').config()
require('isomorphic-fetch')
require('newrelic')
const Logger = require('./logger')
const bugsnag = require('bugsnag')
const db = require('./db')
const listeners = require('./listeners')

const startApp = async () => {
  bugsnag.register(process.env.BUGSNAG_KEY)
  Logger.info('init database')
  await db.init()
  Logger.info('start listeners')
  await listeners.start()
}

startApp()
