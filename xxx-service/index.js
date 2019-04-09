const mesg = require('mesg-js').service()

mesg.listenTask({
  taskX: require('./tasks/taskX')
})
  .on('error', (error) => console.error(error))

mesg.emitEvent('started', { x: true })
  .catch((error) => console.error(error))

console.log('gnarf gnarf')
