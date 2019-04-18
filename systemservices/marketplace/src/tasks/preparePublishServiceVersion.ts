import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { stringToHex, CreateTransaction, hashToHex } from "../contracts/utils";
import { Manifest } from "../types/manifest";

export default (
  marketplace: Marketplace,
  createTransaction: CreateTransaction
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    const manifestHash = await uploadManifest({
      version: '1',
      service: inputs.service
    })
    const transactionData = marketplace.methods.publishServiceVersion(
      stringToHex(inputs.service.definition.sid),
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
  const IPFS = ipfsClient('ipfs.app.mesg.com', '5001', {protocol: 'http'})
  const buffer = Buffer.from(JSON.stringify(manifest))
  const res = await IPFS.add(buffer, {pin: false})
  if (!res.length) {
    throw new Error('Error while uploading manifest')
  }
  return res[0].hash
}
