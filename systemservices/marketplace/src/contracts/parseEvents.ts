import { hexToString, fromUnit, parseTimestamp, hexToHash } from "./utils";

const serviceCreated = (data: any) => ({
  sid: hexToString(data.sid),
  owner: data.owner,
})
const serviceOfferCreated = (data: any) => ({
    sid: hexToString(data.sid),
    offerIndex: data.offerIndex,
    price: fromUnit(data.price),
    duration: data.duration,
})
const serviceOfferDisabled = (data: any) => ({
    sid: hexToString(data.sid),
    offerIndex: data.offerIndex,
})
const serviceOwnershipTransferred = (data: any) => ({
    sid: hexToString(data.sid),
    previousOwner: data.previousOwner,
    newOwner: data.newOwner,
})
const servicePurchased = (data: any) => ({
    sid: hexToString(data.sid),
    offerIndex: data.offerIndex,
    purchaser: data.purchaser,
    price: fromUnit(data.price),
    duration: data.duration,
    expire: parseTimestamp(data.expire),
})
const serviceVersionCreated = (data: any) => ({
    sid: hexToString(data.sid),
    versionHash: hexToHash(data.versionHash),
    manifest: hexToString(data.manifest),
    manifestProtocol: hexToString(data.manifestProtocol),
})

export {
  serviceCreated,
  serviceOfferCreated,
  serviceOfferDisabled,
  serviceOwnershipTransferred,
  servicePurchased,
  serviceVersionCreated,
}
