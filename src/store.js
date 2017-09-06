const utils = require('./utils')
const store = {}

module.exports = web3 => {
  const remove = triggerId => {
    const listener = store[triggerId]
    if (listener) {
      try {
        listener.stopWatching()
      } catch (e) { }
    }
    delete store[triggerId]
  }

  const add = (trigger, onEvent) => {
    const contract = utils.getContract(web3, trigger)
    if (!utils.isValidEvent(contract, trigger.eventName)) { return }  
    remove(trigger.id)
    store[trigger.id] = utils
      .createListenerForEvent(contract, trigger.eventName)

    store[trigger.id]
      .watch((err, event) => onEvent(err, {
        event,
        trigger
      }))
  }

  return {
    add,
    remove
  }
}