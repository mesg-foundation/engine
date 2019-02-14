import Service, { EmitEventReply } from "mesg-js/lib/service/service"
import { EventLog } from "web3/types";
import { hexToAscii } from "../contracts/utils";

export = (mesg: Service, event: EventLog): Promise<EmitEventReply | Error> => {
  return mesg.emitEvent('serviceCreated', {
    sid: hexToAscii(event.returnValues.sid),
    sidHash: event.returnValues.sidHash,
    owner: event.returnValues.owner,
  })
}
