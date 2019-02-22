import Web3 from "web3"
import BigNumber from "bignumber.js";

BigNumber.config({ EXPONENTIAL_AT: 100 })

const maxUint256 = new BigNumber("3963877391197344453575983046348115674221700746820753546331534351508065746944")

const hexToAscii = (x: string) => {
  if (!x) return ""
  return Web3.utils.hexToAscii(x).replace(/\u0000/g, '')
}

const asciiToHex = (x: string) => Web3.utils.asciiToHex(x)
const sha3 = (x: string) => Web3.utils.sha3(x)
const isValidNumber = (x: BigNumber) => !x.isEqualTo(maxUint256)

const toUnit = (x: string|BigNumber) => {
  const n = new BigNumber(x).times(1e18)
  if (!n.integerValue().eq(n)) {
    throw new Error('Number of decimals of ' + x + ' is higher than 18')
  }
  return n
}
const fromUnit = (x: string|BigNumber) => new BigNumber(x).dividedBy(1e18)

const parseTimestamp = (x: string) => new Date(new BigNumber(x).times(1000).toString())

export {
  hexToAscii,
  asciiToHex,
  sha3,
  isValidNumber,
  toUnit,
  fromUnit,
  parseTimestamp,
}
