const { notify } = require('bugsnag')

const NoListenersError = message => {
  this.name = 'NoListenersError'
  this.message = (message || '')
  notify(this.name, this.message)
}
NoListenersError.prototype = Error.prototype

module.exports = NoListenersError
