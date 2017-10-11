const { notify } = require('bugsnag')

function ConnectionError (message) {
  this.name = 'ConnectionError'
  this.message = (message || '')
  notify(this.name, this.message)
}
ConnectionError.prototype = Error.prototype

module.exports = ConnectionError
