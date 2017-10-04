const { notify } = require('bugsnag')

function InvalidBlockchainError (message) {
  this.name = 'InvalidBlockchainError'
  this.message = (message || '')
  notify(this.name, this.message)
}
InvalidBlockchainError.prototype = Error.prototype

module.exports = InvalidBlockchainError
