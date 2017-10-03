const { NoListenersError } = require('../errors')
const listenersModules = [
  require('./ethereum/contractEvent'),
  require('./ethereum/transaction')
]

module.exports = async trigger => {
  const matchingListenersPromises = listenersModules
    .filter(x => x.match(trigger))
    .map(x => x.createListener(trigger))
  
  if (!matchingListenersPromises.length) {
    throw new NoListenersError(trigger.id)
  }

  const listeners = await Promise.all(matchingListenersPromises)
  return {
    watch: callback => listeners.map(x => x.watch(callback)),
    stopWatching: () => listeners.map(x => x.stopWatching())
  }
}
