const MESG = require("mesg-js").service();
const sendgrid = require("@sendgrid/mail");

const send = ({ from, to, subject, text }, { success, failure }) => {
  sendgrid.setApiKey("__CHANGE_WITH_YOUR_SENDGRID_API_KEY__");
  sendgrid.send({ from, to, subject, text })
    .then(([response, _]) => success({ status: response.statusCode }))
    .catch((e) => failure({ message: e.toString() }));
}

MESG.listenTask({ send });