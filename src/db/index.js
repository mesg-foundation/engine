const subscribeToDeletion = require('./subscribeToDeletion')
const subscribeToUpdate = require('./subscribeToUpdate')
const fetchAll = require('./fetchAll')
const createEvent = require('./createEvent')

const init = async () => {
  subscribeToUpdate()
  subscribeToDeletion()
  await fetchAll()
}

module.exports = {
  init,
  createEvent
}
