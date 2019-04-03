import { hexToAscii, fromUnit, parseTimestamp } from "./utils";

const serviceCreated = (data: any) => ({
  sid: hexToAscii(data.sid),
  owner: data.owner,
})
const serviceOfferCreated = (data: any) => ({
    sid: hexToAscii(data.sid),
    offerIndex: data.offerIndex,
    price: fromUnit(data.price),
    duration: data.duration,
})
const serviceOfferDisabled = (data: any) => ({
    sid: hexToAscii(data.sid),
    offerIndex: data.offerIndex,
})
const serviceOwnershipTransferred = (data: any) => ({
    sid: hexToAscii(data.sid),
    previousOwner: data.previousOwner,
    newOwner: data.newOwner,
})
const servicePurchased = (data: any) => ({
    sid: hexToAscii(data.sid),
    offerIndex: data.offerIndex,
    purchaser: data.purchaser,
    price: fromUnit(data.price),
    duration: data.duration,
    expire: parseTimestamp(data.expire),
})
const serviceVersionCreated = (data: any) => ({
    sid: hexToAscii(data.sid),
    versionHash: data.versionHash,
    manifest: hexToAscii(data.manifest),
    manifestProtocol: hexToAscii(data.manifestProtocol),
})

export {
  serviceCreated,
  serviceOfferCreated,
  serviceOfferDisabled,
  serviceOwnershipTransferred,
  servicePurchased,
  serviceVersionCreated,
}
