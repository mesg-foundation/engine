module.exports = async () => [
  await require('./ethereum')('MAINNET'),
  await require('./ethereum')('KOVAN')
]
