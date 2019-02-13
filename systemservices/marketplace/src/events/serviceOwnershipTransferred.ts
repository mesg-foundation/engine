import Service, { EmitEventReply } from "mesg-js/lib/service/service"
import { EventLog } from "web3/types";

export = (mesg: Service, event: EventLog): Promise<EmitEventReply | Error> => {
  return mesg.emitEvent('serviceOwnershipTransferred', {
    hashedSid: event.returnValues.hashedSid,
    previousOwner: event.returnValues.previousOwner,
    newOwner: event.returnValues.newOwner,
  })
}
