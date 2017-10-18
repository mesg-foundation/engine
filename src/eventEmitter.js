const { EventEmitter } = require('events')

let eventEmitter = null

const create = () => eventEmitter
  ? null
  : eventEmitter = new EventEmitter()

const emit = (event, data) => {
  if (!eventEmitter) { throw new Error('need to initialize the event emitter with the `create` function') }
  if (!data.type) { throw new Error('emit needs to send the type of network `type` in the data') }
  if (!data.blockchain) { throw new Error('emit needs to send the blockchain `blockchain` in the data') }
  eventEmitter.emit(`${event.toUpperCase()}:${data.type.toUpperCase()}:${data.blockchain.toUpperCase()}`, data)
  eventEmitter.emit(`${event.toUpperCase()}:${data.type.toUpperCase()}`, data)
  eventEmitter.emit(event.toUpperCase(), data)
}

const emitRawBlock = data => emit('RAW_BLOCK', data)
const emitRawTransaction = data => emit('RAW_TRANSACTION', data)

module.exports = {
  create,
  emitRawBlock,
  emitRawTransaction,
  on: (event, callback) => {
    if (!eventEmitter) { throw new Error('need to initialize the event emitter with the `create` function') }  
    eventEmitter.on(event, callback)
  }
}