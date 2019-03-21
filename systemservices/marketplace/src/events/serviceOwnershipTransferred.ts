import Service, { EmitEventReply } from "mesg-js/lib/service/service"
import { EventLog } from "web3/types";
import { hexToAscii } from "../contracts/utils";

export = (mesg: Service, event: EventLog): Promise<EmitEventReply | Error> => {
  return mesg.emitEvent('serviceOwnershipTransferred', {
    sid: hexToAscii(event.returnValues.sid),
    previousOwner: event.returnValues.previousOwner,
    newOwner: event.returnValues.newOwner,
    transactionHash: event.transactionHash,
    blockNumber: event.blockNumber,
  })
}
