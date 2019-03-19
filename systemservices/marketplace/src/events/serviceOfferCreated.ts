import Service, { EmitEventReply } from "mesg-js/lib/service/service"
import { EventLog } from "web3/types";
import { fromUnit, hexToAscii } from "../contracts/utils";

export = (mesg: Service, event: EventLog): Promise<EmitEventReply | Error> => {
  return mesg.emitEvent('serviceOfferCreated', {
    sid: hexToAscii(event.returnValues.sid),
    offerIndex: event.returnValues.offerIndex,
    price: fromUnit(event.returnValues.price),
    duration: event.returnValues.duration,
  })
}
