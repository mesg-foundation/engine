import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { stringToHex, CreateTransaction, fromUnit, toUnit } from "../contracts/utils";
import { ERC20 } from "../contracts/ERC20";
import BigNumber from "bignumber.js";
import { getService } from "../contracts/service";

export default (
  marketplace: Marketplace,
  token: ERC20,
  createTransaction: CreateTransaction
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  const transactions: Promise<any>[] = []
  let shiftNonce = 0
  try {
    // get service
    const service = await getService(marketplace, inputs.sid)
    
    // check ownership
    if (service.owner.toLowerCase() === inputs.from.toLowerCase()) throw new Error(`service's owner cannot purchase its own service`)

    // get offer data
    const offerIndex = new BigNumber(inputs.offerIndex).toNumber()
    if (offerIndex < 0 || offerIndex >= service.offers.length) {
      throw new Error('offer index is out of range')
    }
    const offer = service.offers[offerIndex]

    // check if offer is active
    if (!offer.active) throw new Error('offer is not active')

    // check user balance
    const balance = fromUnit(await token.methods.balanceOf(inputs.from).call())
    if (offer.price.isGreaterThan(balance)) {
      throw new Error(`purchaser does not have enough balance, needs ${offer.price.toString()} MESG Token`)
    }

    // check allowance balance
    const allowance = fromUnit(await token.methods.allowance(inputs.from, marketplace.options.address).call())
    if (offer.price.isGreaterThan(allowance)) {
      // approve marketplace to spend purchaser token
      const approveTransactionData = token.methods.approve(
        marketplace.options.address, 
        toUnit(offer.price)
      ).encodeABI()
      transactions.push(createTransaction(token, inputs, approveTransactionData))
      shiftNonce++
    }

    // purchase
    const purchaseTransactionData = marketplace.methods.purchase(
      stringToHex(inputs.sid),
      inputs.offerIndex
    ).encodeABI()
    transactions.push(createTransaction(marketplace, inputs, purchaseTransactionData, shiftNonce))

    return outputs.success({
      transactions: await Promise.all(transactions)
    })
  }
  catch (error) {
    console.error('error in preparePurchase', error)
    return outputs.error({ message: error.toString() })
  }
}
