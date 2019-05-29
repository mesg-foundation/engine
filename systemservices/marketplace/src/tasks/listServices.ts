import { TaskInputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { getAllServices } from "../contracts/service";

export default (contract: Marketplace) => async (inputs: TaskInputs): Promise<object> => {
  const services = await getAllServices(contract)
  return { services }
}
