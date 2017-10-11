const { notify } = require('bugsnag')

function InvalidClientError (message) {
  this.name = 'InvalidClientError'
  this.message = (message || '')
  notify(this.name, this.message)
}
InvalidClientError.prototype = Error.prototype

module.exports = InvalidClientError
