import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { stringToHex, CreateTransaction, fromUnit, toUnit } from "../contracts/utils";
import { ERC20 } from "../contracts/ERC20";
import BigNumber from "bignumber.js";
import { getService } from "../contracts/service";
import * as assert from "assert";
import { getServiceOffer } from "../contracts/offer";

export default (
  marketplace: Marketplace,
  token: ERC20,
  createTransaction: CreateTransaction
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  const transactions: Promise<any>[] = []
  let shiftNonce = 0
  try {
    // inputs
    const sid = inputs.sid
    const from = inputs.from
    const offerIndex = new BigNumber(inputs.offerIndex)

    // get service
    const service = await getService(marketplace, sid)
    
    // check ownership
    assert.notStrictEqual(from.toLowerCase(), service.owner.toLowerCase(), `service's owner cannot purchase its own service`)

    // get offer data
    const offer = await getServiceOffer(marketplace, sid, offerIndex)

    // check if offer is active
    assert.ok(offer.active, 'offer is not active')

    // check user balance
    const balance = fromUnit(await token.methods.balanceOf(from).call())
    assert.ok(balance.isGreaterThanOrEqualTo(offer.price), `purchaser does not have enough balance, needs ${offer.price.toString()} MESG Token`)

    // check allowance balance
    const allowance = fromUnit(await token.methods.allowance(from, marketplace.options.address).call())
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
    // QUICK FIX: set a default value to inputs.gas to avoid calling web3.eth.estimateGas that failed if the user's wallet didn't approve the marketplace to spend its token.
    inputs.gas = inputs.gas || 200000
    transactions.push(createTransaction(marketplace, inputs, purchaseTransactionData, shiftNonce))

    return outputs.success({
      transactions: await Promise.all(transactions)
    })
  }
  catch (error) {
    console.error('error in preparePurchase', error)
    return outputs.error({ message: error.message })
  }
}
