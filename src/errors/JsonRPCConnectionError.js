const { notify } = require('bugsnag')

function JsonRPCConnectionError (message) {
  this.name = 'JsonRPCConnectionError'
  this.message = (message || '')
  notify(this.name, this.message)
}
JsonRPCConnectionError.prototype = Error.prototype

module.exports = JsonRPCConnectionError
