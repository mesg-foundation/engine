const express = require('express')
const bodyParser = require('body-parser')
const requestId = require('express-request-id')
const rateLimit = require('express-redis-rate-limit')
const redis = require('../../db/redis')
const { emitRawTransaction } = require('../../eventEmitter')
const Logger = require('../../logger')

const handleRequest = app => (req, res) => {
  const body = req.body
  emitRawTransaction({ type: 'HTTP', request: req, app })
  res.json({
    body: req.body,
    params: req.params,
    data: req.data
  })
}

module.exports = async () => {
  const app = express()
  app.use(requestId())
  app.use(rateLimit(redis, {
    requestLimit: 10,
    timeWindow: 30
  }))
  app.use(bodyParser.json())
  app.use(bodyParser.urlencoded({ extended: true }))
  
  app.post('/triggers/:key', handleRequest(app))
  return new Promise((resolve, reject) => app
    .listen(3000, (err, server) => err
      ? reject(err)
      : resolve(server)
    ))
      .then(_ => Logger.info('Server listening'))
}