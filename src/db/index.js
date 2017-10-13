const Logger = require('../logger')
const Store = require('../store')
const client = require('./client')
const queries = {
  fetchAll: require('./queries/fetchAll'),
  onUpdate: require('./queries/updated'),
  onDelete: require('./queries/deleted'),
  createEvent: require('./queries/mutationEvent')
}

const fetchAll = () => client
  .query({ query: queries.fetchAll })
  .then(x => x.data.allTriggers)

const onDataUpdated = callback => client
  .subscribe({ query: queries.onUpdate })
  .subscribe({
    next: value => callback(null, value.Trigger.node),
    error: error => callback(error, null)
  })

const onDataDeleted = callback => client
  .subscribe({ query: queries.onDelete })
  .subscribe({
    next: value => callback(null, value.Trigger.previousValues.id),
    error: error => callback(error, null)
  })

const writeEvent = (trigger, event) => client
  .mutate({
    mutation: queries.createEvent,
    variables: {
      ...event,
      triggerId: trigger.id
    }
  })
  .then(x => x.data.createEvent)

const init = async () => {
  Logger.info('Fetching triggers...')
  const triggers = await fetchAll()
  triggers.map(Store.add)

  Logger.info('Connecting to trigger update...')
  onDataUpdated((err, trigger) => err
  ? Logger.error(err)
  : Store.update(trigger))

  Logger.info('Connecting to trigger delete...')
  onDataDeleted((err, triggerId) => err
    ? Logger.error(err)
    : Store.remove(triggerId))
}

module.exports = {
  init,
  writeEvent
}
