import { hexToAscii, fromUnit, findInAbi, parseTimestamp } from "../contracts/utils";
import { ABIDefinition } from "web3/eth/abi";
import _marketplaceABI from "../contracts/Marketplace.abi.json"
const marketplaceABI = _marketplaceABI as ABIDefinition[]

const eventHandlers: {[ethName: string]: {
  mesgName: string,
  parse: (event: any) => any
  abi: ABIDefinition
}} = {
  'ServiceCreated': {
    mesgName: 'serviceCreated',
    parse: (event: any) => ({
      sid: hexToAscii(event.sid),
      owner: event.owner,
    }),
    abi: findInAbi(marketplaceABI, 'ServiceCreated')
  },
  'ServiceOfferCreated': {
    mesgName: 'serviceOfferCreated',
    parse: (event: any) => ({
      sid: hexToAscii(event.sid),
      offerIndex: event.offerIndex,
      price: fromUnit(event.price),
      duration: event.duration,
    }),
    abi: findInAbi(marketplaceABI, 'ServiceOfferCreated')
  },
  'ServiceOfferDisabled': {
    mesgName: 'serviceOfferDisabled',
    parse: (event: any) => ({
      sid: hexToAscii(event.sid),
      offerIndex: event.offerIndex,
    }),
    abi: findInAbi(marketplaceABI, 'ServiceOfferDisabled')
  },
  'ServiceOwnershipTransferred': {
    mesgName: 'serviceOwnershipTransferred',
    parse: (event: any) => ({
      sid: hexToAscii(event.sid),
      previousOwner: event.previousOwner,
      newOwner: event.newOwner,
    }),
    abi: findInAbi(marketplaceABI, 'ServiceOwnershipTransferred')
  },
  'ServicePurchased': {
    mesgName: 'servicePurchased',
    parse: (event: any) => ({
      sid: hexToAscii(event.sid),
      offerIndex: event.offerIndex,
      purchaser: event.purchaser,
      price: fromUnit(event.price),
      duration: event.duration,
      expire: parseTimestamp(event.expire),
    }),
    abi: findInAbi(marketplaceABI, 'ServicePurchased')
  },
  'ServiceVersionCreated': {
    mesgName: 'serviceVersionCreated',
    parse: (event: any) => ({
      sid: hexToAscii(event.sid),
      versionHash: event.versionHash,
      manifest: hexToAscii(event.manifest),
      manifestProtocol: hexToAscii(event.manifestProtocol),
    }),
    abi: findInAbi(marketplaceABI, 'ServiceVersionCreated')
  },
}

export {
  eventHandlers
}
