import Web3 from "web3"
import BigNumber from "bignumber.js";
import Contract from "web3/eth/contract";
import { TaskInputs } from "mesg-js/lib/service";
import { Tx } from "web3/eth/types";
import { ABIDefinition } from "web3/eth/abi";
const base58 = require('base-x')('123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz')

BigNumber.config({ EXPONENTIAL_AT: 100 })

const hexToAscii = (x: string) => {
  if (!x) return ''
  return Web3.utils.hexToAscii(x).replace(/\u0000/g, '')
}

const asciiToHex = (x: string) => Web3.utils.asciiToHex(x)
const sha3 = (x: string) => Web3.utils.sha3(x)

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
) => {
  return {
    chainID: chainID,
    to: contract.options.address,
    nonce: (await web3.eth.getTransactionCount(inputs.from)) + (shiftNonce || 0),
    gas: inputs.gas || defaultGas,
    gasPrice: inputs.gasPrice || defaultGasPrice,
    value: '0',
    data: data
  }
}

const findInAbi = (abi: ABIDefinition[], name: string): ABIDefinition => {
  const filter = abi.filter(a => a.name === name)
  if (filter.length !== 1) throw new Error('Did not find definition "'+name+'" in abi')
  return filter[0]
}

export {
  hexToAscii,
  asciiToHex,
  sha3,
  toUnit,
  fromUnit,
  parseTimestamp,
  createTransactionTemplate,
  CreateTransaction,
  hashToHex,
  hexToHash,
  findInAbi,
}
