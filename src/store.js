const store = []

const remove = triggerId => {
  const i = store.findIndex(x => x.id === triggerId)
  if (i >= 0) store.slice(i, 1)
}

const add = trigger => {
  const i = store.findIndex(x => x.id === trigger.id)
  i >= 0 ? store[i] = trigger : store.push(trigger)
  return trigger
}

const update = trigger => trigger.enable ? add(trigger) : remove(trigger.id)

module.exports = ({
  add,
  remove,
  update
})
