import Web3 from 'web3'
import BigNumber from 'bignumber.js';
import Contract from 'web3/eth/contract';
import { TaskInputs } from 'mesg-js/lib/service';
import { Tx } from 'web3/eth/types';
import { ABIDefinition } from 'web3/eth/abi';
import { Log } from 'web3/types';
import * as assert from "assert";
const base58 = require('base-x')('123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz')

BigNumber.config({ EXPONENTIAL_AT: 100 })

const hexToString = Web3.utils.hexToString
const stringToHex = Web3.utils.stringToHex
const keccak256 = Web3.utils.soliditySha3

const toUnit = (x: string|BigNumber): string => {
  const n = new BigNumber(x).times(1e18)
  assert.ok(n.integerValue().eq(n), 'number cannot have more than 18 decimals')
  assert.ok(n.isPositive() || n.isZero(), 'number must be positive or null')
  return n.toString()
}
const fromUnit = (x: string|BigNumber) => new BigNumber(x).dividedBy(1e18)

const parseTimestamp = (x: string) => new Date(new BigNumber(x).times(1000).toNumber())

const hashToHex = (x: string): string => {
  assert.ok(!x.startsWith('0x'), 'hash format is invalid, it starts with 0x')
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
  defaultGasPrice: number
): CreateTransaction => async (
  contract: Contract,
  inputs: TaskInputs,
  data: string,
  shiftNonce?: number
) => {
  const tx = {
    chainID: chainID,
    to: contract.options.address,
    value: '0',
    data: data,
    gas: inputs.gas,
    gasPrice: inputs.gasPrice || defaultGasPrice,
    nonce: (await web3.eth.getTransactionCount(inputs.from)) + (shiftNonce || 0),
  }
  if (!tx.gas) {
    tx.gas = await web3.eth.estimateGas({...tx, from: inputs.from})
  }
  return tx
}

const extractEventFromLogs = (web3: Web3, contract: Contract, eventName: string, logs: Log[]): any => {
  const abi = findInAbi(contract.options.jsonInterface, eventName)
  const log = findInLogs(web3, abi, logs)
  return decodeLog(web3, abi, log)
}

const findInLogs = (web3: Web3, abi: ABIDefinition, logs: Log[]) => {
  const eventSignature = web3.eth.abi.encodeEventSignature(abi)
  const index = logs.findIndex(log => log.topics[0] === eventSignature)
  assert.notStrictEqual(index, -1, `event '${abi.name}' not found in logs`)
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
  assert.notStrictEqual(index, -1, `definition '${eventName}' not found in abi`)
  return abi[index]
}

export {
  hexToString,
  stringToHex,
  keccak256,
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
