const createClient = require('./client')

const listeners = {
  // chain: {
  //   listener: ...,
  //   addresses: {
  //     '0x...': {
  //       address: '0x...',
  //       balance: xxx,
  //       callbacks: [
  //         callback1,
  //         callback2
  //       ]
  //     }
  //   }
  // }
}

const ensureChain = chain => {
  if (listeners[chain]) { return }
  listeners[chain] = {
    client: null,
    addresses: {}
  }
}

const checkAddressBalance = (client, transaction, chain, address) => {
  const { balance, callbacks } = listeners[chain].addresses[address]
  const newBalance = client.eth.getBalance(address)
  if (newBalance.equals(balance)) { return }
  listeners[chain].addresses[address].balance = newBalance
  callbacks.forEach(callback => client
    .handleEvent(callback)(null, {
      ...transaction,
      args: {
        balance: {
          wei: newBalance,
          ether: client.fromWei(newBalance, 'ether')
        }
      }
    })
  )
}

const startListener = async (chain, client) => {
  if (listeners[chain].listener) { return }
  listeners[chain].listener = client.eth.filter({
    fromBlock: 'latest',
    toBlock: 'latest'
  })
  listeners[chain].listener.watch((error, transaction) => {
    if (error) { throw new Error('TODO') }
    Object.keys(listeners[chain].addresses)
      .forEach(address => checkAddressBalance(client, transaction, chain, address))
  })
}

const subscribeToBalanceChange = (address, chain) => async callback => {
  const client = await createClient(chain)
  ensureChain(chain)
  if (!listeners[chain].addresses[address]) {
    listeners[chain].addresses[address] = {
      address,
      balance: client.eth.getBalance(address),
      callbacks: []
    }
  }
  listeners[chain].addresses[address].callbacks.push(callback)
  await startListener(chain, client)
}

const unsubscribeToBalanceChange = (address, chain) => {
  // TODO
}

const createListener = async trigger => {
  const { address, chain } = trigger.connector.ethereumTransaction

  return {
    watch: subscribeToBalanceChange(address, chain),
    stopWatching: unsubscribeToBalanceChange(address, chain)
  }
}

module.exports = {
  match: trigger => trigger.connector.connectorType === 'ETHEREUM_TRANSACTION',
  createListener
}
