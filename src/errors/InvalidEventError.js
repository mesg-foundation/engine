const { notify } = require('bugsnag')

const InvalidEventError = message => {
  this.name = 'InvalidEventError'
  this.message = (message || '')
  notify(this.name, this.message)
}
InvalidEventError.prototype = Error.prototype

module.exports = InvalidEventError
