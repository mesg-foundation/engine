const MESG = require('mesg-js').service()
const Web3 = require('web3')
const web3 = new Web3('wss://mainnet.infura.io/ws')
const contract = new web3.eth.Contract(require('./erc20-abi.json'), '0xe41d2489571d322189246dafa5ebde1f4699f498')

contract.events.Transfer({fromBlock: 'latest'})
.on('data', (event) => {
  console.log('New ERC20 transfer received with hash:', event.transactionHash)
  MESG.emitEvent('transfer', {
    blockNumber: event.blockNumber,
    transactionHash: event.transactionHash,
    from: event.returnValues.from,
    to: event.returnValues.to,
    value: String(event.returnValues.value / Math.pow(10, 18)) // We convert value to its user representation based on the number of decimals used by this ERC20.
  })
})
console.log('Listening ERC20 transfer...')
