const MESG = require('mesg-js').application()

// Event we need to listen
const erc20Transfer = {
  serviceID: '__ERC20_SERVICE_ID__', // The serviceID of the ERC20 service deployed
  eventKey: 'transfer' // The event we want to listen
}

// Task to execute
const sendEmail = {
  serviceID: '__SENDGRID_SERVICE_ID__', // The serviceID of the service to send emails
  taskKey: 'send', // The task we want to execute
  inputs: (eventKey, eventData) => { // This function returns the inputs for of task send based on the data of the event
    console.log('New ERC20 transfer received. will send an email. Transaction hash:', eventData.transactionHash)
    return {
      apiKey: '__SENDGRID_API_KEY__',
      from: 'test@erc20notification.com',
      to: '__REPLACE_WITH_YOUR_EMAIL__',
      subject: 'New ERC20 transfer',
      text: `Transfer from ${eventData.from} to ${eventData.to} of ${eventData.value} tokens -> ${eventData.transactionHash}`
    }
  }
}

MESG.whenEvent(erc20Transfer, sendEmail)
console.log('Listening ERC20 transfer...')
