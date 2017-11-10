const Block = require('../db/block')

module.exports = async (key, fetch, process) => {
  let current = await fetch()
  let last = await Block.last(key)
  if (!last) { return }
  while (last + 1 <= current) {
    await process(last + 1)
    current = await fetch()
    last = await Block.last(key)
  }
}
