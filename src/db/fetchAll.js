const Logger = require('../logger')
const Store = require('../store')
const client = require('./client')
const gql = require('graphql-tag')

const queryCount = gql`query {
  _allTriggersMeta(filter: {
    enable: true
  }) {
    count
  }
}`

const query = gql`query($skip: Int, $first: Int) {
  allTriggers(skip: $skip, first: $first, filter: {
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
  const { count } = (await client.query({ query: queryCount })).data._allTriggersMeta
  const pagination = parseInt(process.env.PAGINATION, 10) || 500
  const pageCount = Math.ceil(count / pagination)
  Logger.info(`Fetching triggers... ${pageCount} pages (${pagination} triggers / page)`)
  const paginationPromise = i => client.query({
    query,
    variables: {
      first: pagination,
      skip: i * pagination
    }
  })
    .then(({ data }) => data.allTriggers.map(Store.add))
  await Promise.all(
    new Array(pageCount)
      .fill()
      .map((_, i) => paginationPromise(i))
  )
}
