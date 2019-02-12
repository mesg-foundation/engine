import Service, { EmitEventReply } from "mesg-js/lib/service/service"
import { EventLog } from "web3/types";

export = (mesg: Service, event: EventLog): Promise<EmitEventReply | Error> => {
  console.log('event', event)
  return mesg.emitEvent('serviceCreated', event.returnValues)
}
