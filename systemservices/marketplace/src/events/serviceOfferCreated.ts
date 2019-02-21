import Service, { EmitEventReply } from "mesg-js/lib/service/service"
import { EventLog } from "web3/types";
import { fromUnit } from "../contracts/utils";

export = (mesg: Service, event: EventLog): Promise<EmitEventReply | Error> => {
  return mesg.emitEvent('serviceOfferCreated', {
    sidHash: event.returnValues.sidHash,
    offerIndex: event.returnValues.offerIndex,
    price: fromUnit(event.returnValues.price),
    duration: event.returnValues.duration,
  })
}
