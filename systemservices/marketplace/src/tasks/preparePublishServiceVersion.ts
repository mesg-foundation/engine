import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { stringToHex, CreateTransaction } from "../contracts/utils";
import { Manifest } from "../types/manifest";
import { getService, isServiceExist } from "../contracts/service";

export default (
  marketplace: Marketplace,
  createTransaction: CreateTransaction
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    // check inputs
    const sid = inputs.service.definition.sid
    if (!(1 <= sid.length() && sid.length() <= 63)) throw new Error('Sid must have a length between 1 and 63')
    // TODO: add check on SID is domain name (see go implementation)

    if(await isServiceExist(marketplace, sid)) {
      // get service
      const service = await getService(marketplace, sid)

      // check ownership
      if (service.owner.toLowerCase() !== inputs.from.toLowerCase()) throw new Error(`Service's owner is different that the specified 'from'`)
    }

    // upload manifest
    const manifestHash = await uploadManifest({
      version: '1',
      service: inputs.service
    })

    // create transaction
    const transactionData = marketplace.methods.publishServiceVersion(
      stringToHex(sid),
      stringToHex(manifestHash),
      stringToHex('ipfs')
    ).encodeABI()
    return outputs.success(await createTransaction(marketplace, inputs, transactionData))
  }
  catch (error) {
    console.error('error in preparePublishServiceVersion', error)
    return outputs.error({ message: error.toString() })
  }
}

const uploadManifest = async (manifest: Manifest): Promise<string> => {
  const ipfsClient = require('ipfs-http-client')
  const IPFS = ipfsClient('ipfs.infura.io', '5001', {protocol: 'https'})
  const buffer = Buffer.from(JSON.stringify(manifest))
  const res = await IPFS.add(buffer, {pin: false})
  if (!res.length) {
    throw new Error('Error while uploading manifest')
  }
  return res[0].hash
}