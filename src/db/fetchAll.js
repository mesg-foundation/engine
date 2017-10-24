const Logger = require('../logger')
const Store = require('../store')
const client = require('./client')
const gql = require('graphql-tag')

const query = gql`query {
  allTriggers(filter: {
    enable: true
  }) {
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
      }
    }
  }
}`

module.exports = async () => {
  Logger.info('Fetching triggers...')
  const { data } = await client.query({ query })
  data.allTriggers.map(Store.add)
}
