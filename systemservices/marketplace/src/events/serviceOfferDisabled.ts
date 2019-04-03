import Service, { EmitEventReply } from "mesg-js/lib/service/service"
import { EventLog } from "web3/types";
import { serviceOfferDisabled } from "../contracts/parseEvents";

export = (mesg: Service, event: EventLog): Promise<EmitEventReply | Error> => {
  return mesg.emitEvent('serviceOfferDisabled', serviceOfferDisabled(event.returnValues))
}
