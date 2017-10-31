const Logger = require('./logger')
const triggerFactory = require('./triggers')

let _store = []

const setStore = store => (_store = store)
const store = () => _store

const remove = triggerId => {
  Logger.info(`Trigger removed`, { triggerId })
  setStore(store().filter(x => x.id !== triggerId))
}

const add = trigger => {
  Logger.info(`Trigger added/replaced`, { triggerId: trigger.id })
  setStore([
    ...store().filter(x => x.id !== trigger.id),
    triggerFactory(trigger)
  ])
  return trigger
}

const update = trigger => trigger.enable
  ? add(trigger)
  : remove(trigger.id)

const all = store

module.exports = {
  add,
  remove,
  update,
  all
}
