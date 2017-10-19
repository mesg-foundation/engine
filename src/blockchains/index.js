module.exports = async () => [
  await require('./ethereum')({ blockchain: 'MAINNET' }),
  await require('./ethereum')({ blockchain: 'TESTNET' })
].filter(x => x)
