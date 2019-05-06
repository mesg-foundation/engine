import Web3 from "web3"
import { EventLog } from "web3/types"
import { Service, EmitEventReply } from "mesg-js/lib/service"

import { newBlockEventEmitter } from "../newBlock"
import { Marketplace } from "../contracts/Marketplace";

import serviceCreated from "./serviceCreated"
import serviceOfferCreated from "./serviceOfferCreated"
import serviceOfferDisabled from "./serviceOfferDisabled"
import serviceOwnershipTransferred from "./serviceOwnershipTransferred"
import serviceVersionCreated from "./serviceVersionCreated"
import servicePurchased from "./servicePurchased"

const eventHandlers: {[eventName: string]: (mesg: Service, event: EventLog) => Promise<EmitEventReply | Error>} = {
  'ServiceCreated': serviceCreated,
  'ServiceOfferCreated': serviceOfferCreated,
  'ServiceOfferDisabled': serviceOfferDisabled,
  'ServiceOwnershipTransferred': serviceOwnershipTransferred,
  'ServicePurchased': servicePurchased,
  'ServiceVersionCreated': serviceVersionCreated,
}

export default async (
  mesg: Service,
  web3: Web3,
  marketplace: Marketplace,
  blockConfirmations: number,
  pollingTime: number
) => {
  const newBlock = await newBlockEventEmitter(web3, blockConfirmations, null, pollingTime)
  newBlock.on('newBlock', async blockNumber => {
    try {
      console.log('new block', blockNumber)
      const events = await marketplace.getPastEvents("allEvents", {
        fromBlock: blockNumber,
        toBlock: blockNumber,
      })
      events.forEach(async event => {
        try {
          if (!eventHandlers[event.event]) {
            throw new Error(`event '${event.event}' not implemented`)
          }
          eventHandlers[event.event](mesg, event)
        } catch(error) {
          console.error(`An error occurred during processing of event '${event.event}'.`, error)
        }
      })
    }
    catch (error) {
      console.error('catch newBlock on', error)
    }
  })
}
