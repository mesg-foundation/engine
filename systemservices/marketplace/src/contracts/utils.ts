import Web3 from "web3"

const hexToAscii = (x: string) => Web3.utils.hexToAscii(x).replace(/\u0000/g, '')
const asciiToHex = (x: string) => Web3.utils.asciiToHex(x)

export {
  hexToAscii,
  asciiToHex,
}
