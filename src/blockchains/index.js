module.exports = async () => [
  await require('./ethereum')({ blockchain: 'MAINNET' }),
  await require('./ethereum')({ blockchain: 'KOVAN' })
].filter(x => x)
