import Service, { EmitEventReply } from "mesg-js/lib/service/service"
import { EventData } from "web3-eth-contract/types";

export = (mesg: Service, event: EventData): Promise<EmitEventReply | Error> => {
  console.log('event', event)
  return mesg.emitEvent('serviceCreated', event.returnValues)
}
