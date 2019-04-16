import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { toUnit, stringToHex, CreateTransaction } from "../contracts/utils";
import BigNumber from "bignumber.js";
import { getService } from "../contracts/service";

export default (
  marketplace: Marketplace,
  createTransaction: CreateTransaction
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    // check inputs
    const duration = new BigNumber(inputs.duration)
    if (duration.isNegative() || duration.isZero()) throw new Error('duration cannot be negative or equal to zero')

    // check service
    const service = await getService(marketplace, inputs.sid)

    // check ownership
    if (service.owner.toLowerCase() !== inputs.from.toLowerCase()) throw new Error(`service's owner is different that the specified 'from'`)

    // check service version
    if (service.versions.length === 0) throw new Error('cannot create an offer on a service with 0 version')

    // create transaction
    const transactionData = marketplace.methods.createServiceOffer(
      stringToHex(inputs.sid),
      toUnit(inputs.price),
      duration.toString()
    ).encodeABI()
    return outputs.success(await createTransaction(marketplace, inputs, transactionData))
  }
  catch (error) {
    console.error('error in prepareCreateServiceOffer', error)
    return outputs.error({ message: error.toString() })
  }
}
