const MESG = require('mesg-js').service()
const sendgrid = require('@sendgrid/mail')

const sendHandler = async (inputs, outputs) => {
  try {
    console.log('New send task received')
    // Configure SendGrid
    sendgrid.setApiKey('__CHANGE_WITH_YOUR_SENDGRID_API_KEY__')
    // Sends an email with the inputs
    const result = await sendgrid.send({
      from: inputs.from,
      to: inputs.to,
      subject: inputs.subject,
      text: inputs.text
    })
    // Return the success output
    return outputs.success({
      status: result[0].statusCode
    })
  } catch (error) {
    // If an error occurs, return the failure output
    return outputs.failure({
      message: error.toString()
    })
  }
}

MESG.listenTask({
  send: sendHandler
})
console.log('Listening tasks...')
