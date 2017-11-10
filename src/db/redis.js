const redis = require('redis')
const Logger = require('../logger')

const client = redis.createClient({ host: 'redis' })
client.on('error', Logger.error)

module.exports = client
