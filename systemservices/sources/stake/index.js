const MESG = require('mesg-js').service()
const Web3 = require('web3')
const erc20ABI = require('./erc20-abi.json')
const web3 = new Web3("https://mainnet.infura.io/")

const balanceOf = async (inputs, outputs) => {
  const contract = new web3.eth.Contract(erc20ABI, inputs.contractAddress)
  return contract.methods.balanceOf(inputs.address).call()
    .then(balance => outputs.success({ balance }))
    .catch(err => outputs.error({ message: err.toString() }))
}

MESG.listenTask({
  balanceOf
})