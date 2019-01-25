import Web3 from "web3"

const hexToAscii = (x: any) => Web3.utils.hexToAscii(x).replace(/\u0000/g, '')

export { hexToAscii }