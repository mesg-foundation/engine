import * as utils from 'web3-utils';

const hexToAscii = (x: string) => utils.hexToAscii(x).replace(/\u0000/g, '')
const asciiToHex = (x: string) => utils.asciiToHex(x)

export {
  hexToAscii,
  asciiToHex,
}
