const { notify } = require('bugsnag')

function InvalidTriggerError (message) {
  this.name = 'InvalidTriggerError'
  this.message = (message || '')
  notify(this.name, this.message)
}
InvalidTriggerError.prototype = Error.prototype

module.exports = InvalidTriggerError
