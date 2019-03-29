import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { asciiToHex, CreateTransaction, hashToHex } from "../contracts/utils";
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
      asciiToHex(inputs.service.definition.sid),
      asciiToHex(manifestHash),
      asciiToHex('ipfs')
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