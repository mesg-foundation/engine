import Service, { EmitEventReply } from "mesg-js/lib/service/service"
import { EventLog } from "web3/types";

export = (mesg: Service, event: EventLog): Promise<EmitEventReply | Error> => {
  return mesg.emitEvent('serviceOfferCreated', {
    sid: event.returnValues.sid,
    offerIndex: event.returnValues.offerIndex,
    price: event.returnValues.price,
    duration: event.returnValues.duration,
  })
}
