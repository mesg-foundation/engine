const MESG = require('mesg-js').service()
const Web3 = require('web3')
const web3 = new Web3('wss://mainnet.infura.io/_ws')
const contract = new web3.eth.Contract(require('./erc20-abi.json'), "0xf230b790e05390fc8295f4d3f60332c93bed42e2")

contract.events.Transfer({fromBlock: 'latest'})
.on('data', event => {
  MESG.emitEvent('transfer', {
    blockNumber: event.blockNumber,
    transactionHash: event.transactionHash,
    from: event.returnValues.from,
    to: event.returnValues.to,
    value: event.returnValues.value / Math.pow(10, 6),
  })
})
