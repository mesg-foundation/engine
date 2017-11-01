const blocks = {
  // [JSON.stringify({ type: 'ETHEREUM', blockchain: 'TESTNET' })]: 4555546
}

const processed = async (key, block) => {
  blocks[JSON.stringify(key)] = block
  console.log(blocks)
}

const last = async key => blocks[JSON.stringify(key)] || null

module.exports = {
  processed,
  last
}
