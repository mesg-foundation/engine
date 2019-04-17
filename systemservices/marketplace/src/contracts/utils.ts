import Web3 from 'web3'
import BigNumber from 'bignumber.js';
import Contract from 'web3/eth/contract';
import { TaskInputs } from 'mesg-js/lib/service';
import { Tx } from 'web3/eth/types';
import { ABIDefinition } from 'web3/eth/abi';
import { Log } from 'web3/types';
const base58 = require('base-x')('123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz')

BigNumber.config({ EXPONENTIAL_AT: 100 })

const hexToString = Web3.utils.hexToString
const stringToHex = Web3.utils.stringToHex
const sha3 = Web3.utils.sha3

const toUnit = (x: string|BigNumber): string => {
  const n = new BigNumber(x).times(1e18)
  if (!n.integerValue().eq(n)) {
    throw new Error('Number of decimals of ' + x + ' is higher than 18')
  }
  return n.toString()
}
const fromUnit = (x: string|BigNumber) => new BigNumber(x).dividedBy(1e18)

const parseTimestamp = (x: string) => new Date(new BigNumber(x).times(1000).toNumber())

const hashToHex = (x: string): string => {
  if (x.startsWith('0x')) throw new Error('hash format is invalid. It starts with 0x')
  return '0x' + base58.decode(x).toString('hex')
}
const hexToHash = (x: string): string => base58.encode(Buffer.from(x.replace(/^0x/, ''), 'hex'))

interface CreateTransaction {
  (
    contract: Contract,
    inputs: TaskInputs,
    data: string,
    shiftNonce?: number
  ): Promise<Tx>;
};

const createTransactionTemplate = (
  chainID: number,
  web3: Web3,
  defaultGas: number,
  defaultGasPrice: number
): CreateTransaction => async (
  contract: Contract,
  inputs: TaskInputs,
  data: string,
  shiftNonce?: number
) => ({
  chainID: chainID,
  to: contract.options.address,
  nonce: (await web3.eth.getTransactionCount(inputs.from)) + (shiftNonce || 0),
  gas: inputs.gas || defaultGas,
  gasPrice: inputs.gasPrice || defaultGasPrice,
  value: '0',
  data: data
})

const extractEventFromLogs = (web3: Web3, contract: Contract, eventName: string, logs: Log[]): any => {
  const abi = findInAbi(contract.options.jsonInterface, eventName)
  const log = findInLogs(web3, abi, logs)
  return decodeLog(web3, abi, log)
}

const findInLogs = (web3: Web3, abi: ABIDefinition, logs: Log[]) => {
  const eventSignature = web3.eth.abi.encodeEventSignature(abi)
  const index = logs.findIndex(log => log.topics[0] === eventSignature)
  if (index === -1) throw new Error(`Did not find event '${abi.name}' in logs`)
  return logs[index]
}

const decodeLog = (web3: Web3, abi: ABIDefinition, log: Log): any => {
  // Remove first element because event is non-anonymous
  // https://web3js.readthedocs.io/en/1.0/web3-eth-abi.html#decodelog
  if (abi.anonymous === false) log.topics.splice(0, 1)
  return web3.eth.abi.decodeLog(abi.inputs as object, log.data, log.topics)
}

const findInAbi = (abi: ABIDefinition[], eventName: string): ABIDefinition => {
  const index = abi.findIndex(a => a.name === eventName)
  if (index === -1) throw new Error(`Did not find definition '${eventName}' in abi`)
  return abi[index]
}

export {
  hexToString,
  stringToHex,
  sha3,
  toUnit,
  fromUnit,
  parseTimestamp,
  createTransactionTemplate,
  CreateTransaction,
  hashToHex,
  hexToHash,
  decodeLog,
  findInAbi,
  findInLogs,
  extractEventFromLogs,
}
