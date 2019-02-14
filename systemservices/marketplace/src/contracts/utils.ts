import Web3 from "web3"

const hexToAscii = (x: string) => Web3.utils.hexToAscii(x).replace(/\u0000/g, '')
const asciiToHex = (x: string) => Web3.utils.asciiToHex(x)
const sha3 = (x: string) => Web3.utils.sha3(x)

export {
  hexToAscii,
  asciiToHex,
  sha3,
}
