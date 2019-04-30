import { TaskInputs, TaskOutputs } from "mesg-js/lib/service"
import { Marketplace } from "../contracts/Marketplace"
import { stringToHex, CreateTransaction } from "../contracts/utils";
import { Manifest } from "../types/manifest";
import { getService, isServiceExist } from "../contracts/service";
import * as assert from "assert";
import { computeVersionHash, isVersionExist } from "../contracts/version";
import isDomainName from "../contracts/isDomainName";

const manifestProtocol = 'ipfs'

export default (
  marketplace: Marketplace,
  createTransaction: CreateTransaction
) => async (inputs: TaskInputs, outputs: TaskOutputs): Promise<void> => {
  try {
    // check inputs
    const sid = inputs.service.definition.sid
    assert.ok(sid.length >= 1 && sid.length <= 63, 'sid\'s length must be 1 at min and 63 at max') // See Core service validation (https://github.com/mesg-foundation/core)
    assert.ok(isDomainName(sid), 'sid must respect domain-style notation, eg author.name')

    if(await isServiceExist(marketplace, sid)) {
      // get service
      const service = await getService(marketplace, sid)

      // check ownership
      assert.strictEqual(inputs.from.toLowerCase(), service.owner.toLowerCase(), `service's owner is different`)
    }

    // upload manifest
    const manifestHash = await uploadManifest({
      version: '1',
      service: inputs.service
    })

    // check if version already exist on marketplace
    const versionHash = computeVersionHash(inputs.from, sid, manifestHash, manifestProtocol)
    assert.ok(!(await isVersionExist(marketplace, versionHash)), `service version already exist`)

    // create transaction
    const transactionData = marketplace.methods.publishServiceVersion(
      stringToHex(sid),
      stringToHex(manifestHash),
      stringToHex(manifestProtocol)
    ).encodeABI()
    return outputs.success(await createTransaction(marketplace, inputs, transactionData))
  }
  catch (error) {
    console.error('error in preparePublishServiceVersion', error)
    return outputs.error({ message: error.message })
  }
}

const uploadManifest = async (manifest: Manifest): Promise<string> => {
  const ipfsClient = require('ipfs-http-client')
  const IPFS = ipfsClient(process.env.IPFS_PROVIDER)
  const buffer = Buffer.from(JSON.stringify(manifest))
  const res = await IPFS.add(buffer, {pin: false})
  if (!res.length) {
    throw new Error('error while uploading manifest')
  }
  return res[0].hash
}
