import Service, { EmitEventReply } from "mesg-js/lib/service/service"
import { EventLog } from "web3/types";

export = (mesg: Service, event: EventLog): Promise<EmitEventReply | Error> => {
  return mesg.emitEvent('serviceOfferDisabled', {
    hashedSid: event.returnValues.hashedSid,
    offerIndex: event.returnValues.offerIndex,
  })
}
