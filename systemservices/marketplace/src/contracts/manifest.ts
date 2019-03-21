import { Validator } from "jsonschema"
import Request from "request-promise-native"
import { Manifest } from "../types/manifest"
import manifestSchema from '../types/schema/manifest.json'
import serviceSchema from '../types/schema/definition.json'

const validator = new Validator();

const getIpfs = async (source: string): Promise<any> => {
  return await Request.get('https://gateway.ipfs.io/ipfs/' + source, { json: true, timeout: 10000 })
}

const getHttp = async (source: string): Promise<any> => {
  return await Request.get(source, { json: true, timeout: 10000 })
}

const get: {[key: string]: (source: string) => Promise<any>} = {
  'ipfs': getIpfs,
  'https': getHttp,
  'http': getHttp,
}

const getManifest = async (protocol: string, source: string): Promise<Manifest|undefined> => {
  try {
    if (!get[protocol]) {
      console.warn('protocol', protocol, 'is not compatible with this service')
      return
    }
    const manifest: Manifest = await get[protocol](source)
    const validation = validator.validate(manifest, manifestSchema)
    if (!validation.valid) {
      console.warn('manifest', protocol, '::', source, 'is not a valid manifest')
      return
    }
    const defValidation = validator.validate(manifest.service.definition, serviceSchema)
    if (!defValidation.valid) {
      console.warn('manifest', protocol, '::', source, 'doesn\'t have a valid service definition')
      return
    }
    return manifest
  }
  catch (error) {
    console.warn('error while downloading manifest', protocol, '::', source, error.toString())
    return
  }
}

export { getManifest }
