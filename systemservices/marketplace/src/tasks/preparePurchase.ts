import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { asciiToHex, CreateTransaction, fromUnit, toUnit } from "../contracts/utils";
import { ERC20 } from "../contracts/ERC20";
import BigNumber from "bignumber.js";
import { getServiceOffer } from "../contracts/offer";

export default (
  marketplace: Marketplace,
  token: ERC20,
  createTransaction: CreateTransaction
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  const transactions: Promise<any>[] = []
  let shiftNonce = 0
  try {
    // get offer data
    const offer = await getServiceOffer(marketplace, inputs.sid, new BigNumber(inputs.offerIndex))
    if (offer === undefined) {
      throw new Error('Offer with index ' + inputs.offerIndex + ' and sid ' + inputs.sid + ' does not exist')
    }

    // check user balance
    const balance = fromUnit(await token.methods.balanceOf(inputs.from).call())
    if (offer.price.isGreaterThan(balance)) {
      throw new Error('Purchaser does not have enough balance')
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
      asciiToHex(inputs.sid),
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
