const winston = require('winston')
require('winston-loggly-bulk')

const logger = new winston.Logger()

if (process.env.NODE_ENV !== 'production') {
  logger.add(winston.transports.Console, {
    timestamp: true,
    colorize: true,
    prettyPrint: true
  })
}

if (process.env.LOGGLY_TOKEN) {
  logger.add(winston.transports.Loggly, {
    token: process.env.LOGGLY_TOKEN,
    subdomain: process.env.LOGGLY_DOMAIN,
    json: true
  })
}

module.exports = logger
