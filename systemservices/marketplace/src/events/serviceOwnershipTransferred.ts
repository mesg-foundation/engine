import Service from "mesg-js/lib/service/service"
import { EventLog } from "web3/types";
import { serviceOwnershipTransferred } from "../contracts/parseEvents";
import { EventCreateOutputs } from "mesg-js/lib/api";

export = (mesg: Service, event: EventLog): EventCreateOutputs => {
  return mesg.emitEvent('serviceOwnershipTransferred', serviceOwnershipTransferred(event.returnValues))
}
