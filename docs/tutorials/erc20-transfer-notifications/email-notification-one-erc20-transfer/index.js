const MESG = require("mesg-js").application();

// Event we need to listen
const erc20Transfer = {
  serviceID: __ERC20_SERVICE_ID__, // The serviceID of the ERC20 service deployed
  filter: "transfer" // The event we want to listen
};

// Task to execute
const sendEmail = {
  serviceID: __SENDGRID_SERVICE_ID__, // The serviceID of the service to send emails
  taskKey: "send", // The task we want to execute
  inputs: (eventKey, { from, to, value, transactionHash }) => { // a function that returns the inputs for the send task based on the data of the event
    console.log("new transfer received with hash", transactionHash);
    return {
      apiKey: __SENDGRID_API_KEY__,
      from: "test@erc20notification.com",
      to: __REPLACE_WITH_YOUR_EMAIL__,
      subject: 'New ERC20 transfer',
      text: `Transfer from ${from} to ${to} of ${value} tokens -> ${transactionHash}`
    };
  }
};

MESG.whenEvent(erc20Transfer, sendEmail);
