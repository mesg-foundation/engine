module.exports = ({ foo, bar }, { success, error }) => {
  if (foo === 'hello' && bar === 'world') {
    success({ message: 'Hello world is valid' })
  } else {
    error({ error: 'invalid inputs' })
  }
}
