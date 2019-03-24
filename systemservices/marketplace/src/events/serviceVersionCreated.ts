import Service, { EmitEventReply } from "mesg-js/lib/service/service"
import { EventLog } from "web3/types";
import { hexToAscii } from "../contracts/utils";
import { getManifest } from "../contracts/manifest";

export = async (mesg: Service, event: EventLog): Promise<EmitEventReply | Error> => {
  const manifest = await getManifest(hexToAscii(event.returnValues.manifestProtocol), hexToAscii(event.returnValues.manifest))
  return mesg.emitEvent('serviceVersionCreated', {
    sid: hexToAscii(event.returnValues.sid),
    versionHash: event.returnValues.versionHash,
    manifest: manifest
  })
}
