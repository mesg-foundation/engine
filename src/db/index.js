const client = require('./client')
const Store = require('../store')
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
  console.debug('Fetching triggers...')
  const triggers = await fetchAll()
  triggers.map(Store.add)
  
  console.debug('Connecting to trigger update...')
  onDataUpdated((err, trigger) => err
  ? console.error(err)
  : Store.update(trigger))
  
  console.debug('Connecting to trigger delete...')
  onDataDeleted((err, triggerId) => err
    ? console.error(err)
    : Store.remove(triggerId))
}

module.exports = {
  init,
  writeEvent
}
