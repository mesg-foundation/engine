require('dotenv').config()
require('isomorphic-fetch')
require('newrelic')
const bugsnag = require('bugsnag')
const db = require('./db')
const listeners = require('./listeners')

const startApp = async () => {
  bugsnag.register(process.env.BUGSNAG_KEY)
  console.debug('init database')
  await db.init()
  console.debug('start listeners')
  await listeners.start()
}

startApp()
