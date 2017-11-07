const client = require('./redis')

const id = ({ type, blockchain }) => [type, blockchain].join('/')

const processed = (key, block) => new Promise((resolve, reject) => client
  .set(id(key), block, (err, data) => err
    ? reject(err)
    : resolve(data)
  ))

const last = key => new Promise((resolve, reject) => client
  .get(id(key), (err, data) => err
    ? reject(err)
    : resolve(parseInt(data, 10) || null)
  ))

module.exports = {
  processed,
  last
}
