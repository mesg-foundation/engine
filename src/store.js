const triggerFactory = require('./triggers')

const store = []

const remove = triggerId => {
  const i = store.findIndex(x => x.id === triggerId)
  if (i >= 0) store.slice(i, 1)
}

const add = trigger => {
  const i = store.findIndex(x => x.id === trigger.id)
  i >= 0 ? store[i] = triggerFactory(trigger) : store.push(triggerFactory(trigger))
  return trigger
}

const update = trigger => trigger.enable ? add(trigger) : remove(trigger.id)

const matchingTriggers = args => store
  .filter(x => x.match(args))

module.exports = ({
  add,
  remove,
  update,
  matchingTriggers
})
