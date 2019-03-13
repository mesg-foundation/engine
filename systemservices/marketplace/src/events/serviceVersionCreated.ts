import Service, { EmitEventReply } from "mesg-js/lib/service/service"
import { EventLog } from "web3/types";
import { hexToAscii } from "../contracts/utils";

export = (mesg: Service, event: EventLog): Promise<EmitEventReply | Error> => {
  return mesg.emitEvent('serviceVersionCreated', {
    sid: hexToAscii(event.returnValues.sid),
    hash: event.returnValues.hash,
    manifest: hexToAscii(event.returnValues.manifest),
    manifestProtocol: hexToAscii(event.returnValues.manifestProtocol),
  })
}
