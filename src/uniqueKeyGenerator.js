const MAX_PAD_NUMBER = 16
const KEY_SEPARATOR = '-'

/**
 * Pad the value with 0 to have a maximum of MAX_PAD_NUMBER digits
 * @param {Number} value 
 * @returns {String} the value as string padded with 0 to have MAX_PAD_NUMBER digits
 */
const padNumber = value => [
  ...new Array(MAX_PAD_NUMBER).fill(0),
  ...(value !== null ? [value] : [])
].slice(-MAX_PAD_NUMBER).join('')

/**
 * Format the value of the input with different process for string or number
 *   - number will be transformed as string and padded with 0 to ensure all will have the same lenght
 *   - string will be returned with no transformations
 * @param {String | Number} value 
 * @returns {String} the formatted value of the value in parameter
 */
const format = value => typeof value === 'string'
  ? value
  : padNumber(value)

/**
 * 
 * @param {Array} keys that will contains the differents elements that will compose a unique key
 *                     all of thoses keys shoud be sortable and the order we give the keys will
 *                     be important for the global sorting
 * @returns {String}
 */
module.exports = keys => keys
  .map(x => format(x))
  .join(KEY_SEPARATOR)

module.exports.KEY_SEPARATOR = KEY_SEPARATOR