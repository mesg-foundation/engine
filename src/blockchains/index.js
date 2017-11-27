module.exports = async () => [
  await require('./ethereum')({ blockchain: 'MAINNET' }),
  await require('./ethereum')({ blockchain: 'TESTNET' }),
  await require('./http/webhook')()
].filter(x => x)
