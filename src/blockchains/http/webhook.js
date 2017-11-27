const express = require('express')
const bodyParser = require('body-parser')
const Logger = require('../../logger')

const handleRequest = (req, res) => {
  const body = req.body
  Logger.info(`request with ${JSON.stringify(req.body)}`)
  Logger.info(`request with ${JSON.stringify(req.params)}`) 
  Logger.info(`request with ${JSON.stringify(req.data)}`) 
  res.json({
    body: req.body,
    params: req.params,
    data: req.data
  })
}

module.exports = async () => {
  const app = express()
  app.use(bodyParser.json())
  app.use(bodyParser.urlencoded({ extended: true }))
  
  app.post('/triggers/:id', handleRequest)
  return new Promise((resolve, reject) => app
    .listen(3000, (err, server) => err
      ? reject(err)
      : resolve(server)
    ))
      .then(_ => Logger.info('Server listening'))
}