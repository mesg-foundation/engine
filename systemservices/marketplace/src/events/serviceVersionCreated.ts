import Service, { EmitEventReply } from "mesg-js/lib/service/service"
import { EventLog } from "web3/types";

export = (mesg: Service, event: EventLog): Promise<EmitEventReply | Error> => {
  return mesg.emitEvent('serviceVersionCreated', {
    sid: event.returnValues.sid, // TODO: to convert to ascii
    hash: event.returnValues.hash,
    manifest: event.returnValues.manifest, // TODO: to convert to ascii
    manifestProtocol: event.returnValues.manifestProtocol, // TODO: to convert to ascii
  })
}
