const listeners = [
  require('./ethereum/contractEvent'),
  require('./ethereum/transaction')
]

module.exports = trigger => {
  const matchingListeners = listeners
    .filter(x => x.match(trigger))
    .map(x => x.createListener(trigger))
  if (!matchingListeners.length) throw new Error(`${trigger.id} does not have any valid listener`)

  return matchingListeners.length === 1
    ? matchingListeners[0]
    : {
      watch: callback => matchingListeners.forEach(x => x.watch(callback)),
      stopWatching: () => matchingListeners.forEach(x => x.stopWatching())
    }
}
