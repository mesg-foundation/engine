const { notify } = require('bugsnag')

const InvalidBlockchainError = message => {
  this.name = 'InvalidBlockchainError'
  this.message = (message || '')
  notify(this.name, this.message)
}
InvalidBlockchainError.prototype = Error.prototype

module.exports = InvalidBlockchainError
