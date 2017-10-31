const Logger = require('../logger')
const subscribeToDeletion = require('./subscribeToDeletion')
const subscribeToUpdate = require('./subscribeToUpdate')
const fetchAll = require('./fetchAll')
const createEvent = require('./createEvent')

const init = async () => {
  try {
    subscribeToUpdate()
    subscribeToDeletion()
    await fetchAll()
  } catch (e) {
    Logger.error('Cannot initialize database')
    throw e
  }
}

module.exports = {
  init,
  createEvent
}
