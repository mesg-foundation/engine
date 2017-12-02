const Logger = require('../logger')
const Store = require('../store')
const { client } = require('./client')
const gql = require('graphql-tag')

const query = id => gql`{
  Trigger(id: "${id}") {
    id
    enable
    connector {
      connectorType
      ethereumContract {
        eventName
        contract {
          abi
          address
          chain
        }
      },
      ethereumToken {
        eventName
        contract {
          abi
          address
          chain
        }
      },
      ethereumTransaction {
        address
        chain
        matchType
      }
      webhook {
        key
      }
    }
  }
}`

module.exports = async id => {
  try {
    Logger.info(`Fetching trigger... ${id}`)
    const { Trigger } = (await client.query({
      query: query(id),
      fetchPolicy: 'network-only'
    })).data
    Trigger.enable
      ? Store.add(Trigger)
      : Store.remove(Trigger.id)
    return Trigger
  } catch (e) {
    Logger.error(`cannot fetch trigger ${id}`)
    throw e
  }
}
