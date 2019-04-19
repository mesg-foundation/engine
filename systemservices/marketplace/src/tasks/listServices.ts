import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { getAllServices } from "../contracts/service";

export default (contract: Marketplace) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const services = await getAllServices(contract)
    return outputs.success({ services })
  }
  catch (error) {
    console.error('error in listServices', error)
    return outputs.error({ message: error.message })
  }
}
