const createClient = require('./client')

let subscribers = [
  // { address: 'xxx', balance: 'xxx', chain: 'xxx', callback: Function, id: 'xx' }
]

let listeners = {
  // KOVAN: xxx, MAINNET: xxx
}

const onNewTransaction = (client, chain) => (error, transaction) => {
  if (error) { throw new Error('TODO') }
  subscribers
    .filter(x => x.chain === chain)
    .filter(updateBalance(client))
    .forEach(notify(client, transaction))
}

const startListener = async (chain, client, onEvent) => {
  if (listeners[chain]) { return }
  listeners[chain] = client.eth.filter({
    fromBlock: 'latest',
    toBlock: 'latest'
  })
  listeners[chain].watch(onNewTransaction(client, chain))
}

const updateBalance = client => subscriber => {
  const balance = client.eth.getBalance(subscriber.address)
  if (subscriber.balance.equals(balance)) { return false }
  Object.assign(subscriber, { balance }) // Object.assign(subscriber... like that it edits directly the value 
  return true
}

const notify = (client, transaction) => subscriber => {
  client.handleEvent(subscriber.callback)(null, {
    ...transaction,
    args: {
      balance: {
        wei: subscriber.balance,
        ether: client.fromWei(subscriber.balance, 'ether')
      }
    }
  })
}

const subscribeToBalanceChange = (address, chain, id) => async callback => {
  const client = await createClient(chain)
  subscribers.push({
    id,
    address,
    chain,
    callback,
    balance: client.eth.getBalance(address)
  })
  await startListener(chain, client)
}

const unsubscribeToBalanceChange = (address, chain, id) => () => {
  subscribers = subscribers.filter(x => x.id !== id)
}

const createListener = async trigger => {
  const { address, chain } = trigger.connector.ethereumTransaction
  return {
    watch: subscribeToBalanceChange(address, chain, trigger.id),
    stopWatching: unsubscribeToBalanceChange(address, chain, trigger.id)
  }
}

module.exports = {
  match: trigger => trigger.connector.connectorType === 'ETHEREUM_TRANSACTION',
  createListener
}
